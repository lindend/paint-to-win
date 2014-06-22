package storage

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

type Storage struct {
	Db    *gorm.DB
	cache *redis.Pool
}

func NewStorage(db *gorm.DB, cacheAddress string) (*Storage, error) {
	cache := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", cacheAddress)
	}, 240)

	return &Storage{
		Db:    db,
		cache: cache,
	}, nil
}

func (storage *Storage) Save(item interface{}) error {
	return storage.Db.Save(item).Error
}

func (storage *Storage) FirstWhere(item interface{}, out interface{}) error {
	return storage.Db.Where(item).First(out).Error
}

func (storage *Storage) Where(item interface{}, out interface{}) error {
	return storage.Db.Where(item).Find(out).Error
}

func (storage *Storage) Exists(item interface{}) (bool, error) {
	var count int
	err := storage.Db.Where(item).Count(&count).Error
	return count > 0, err
}

func (storage *Storage) addToCacheMap(collectionKey string, itemKey string, data interface{}) error {
	cache := storage.cache.Get()
	defer cache.Close()

	var encodedData string
	var err error
	if encodedData, err = storage.encode(data); err != nil {
		return err
	}
	_, err = cache.Do("HSET", collectionKey, itemKey, encodedData)
	return err
}

func (storage *Storage) removeFromCacheMap(collectionKey string, itemKey string) {
	cache := storage.cache.Get()
	defer cache.Close()

	cache.Do("HDEL", collectionKey, itemKey)
}

func (storage *Storage) getCacheMap(collectionKey string) ([]string, error) {
	cache := storage.cache.Get()
	defer cache.Close()

	return redis.Strings(cache.Do("HGETALL"))
}

func (storage *Storage) GetFromCache(key string, val interface{}) error {
	cache := storage.cache.Get()
	defer cache.Close()
	data, err := redis.Bytes(cache.Do("GET", key))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}

func (storage *Storage) SaveInCache(key string, data interface{}, timeout int) error {
	cache := storage.cache.Get()
	defer cache.Close()

	var encoded string
	var err error
	if encoded, err = storage.encode(data); err != nil {
		return err
	}
	_, err = cache.Do("SETEX", key, timeout, encoded)
	return err
}

func (storage *Storage) encode(data interface{}) (string, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
