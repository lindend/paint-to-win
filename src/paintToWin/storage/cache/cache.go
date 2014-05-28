package cache

import (
	"github.com/garyburd/redigo/redis"
)

type cacheStrategy struct {
	timeout int
}

type RedisCache struct {
	cache *redis.Pool
}

func (cache *RedisCache) SetCaching(valueType interface{}) {

}
