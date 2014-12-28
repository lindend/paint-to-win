package service

import (
	"github.com/jinzhu/gorm"

	"paintToWin/storage"
)

type DbServiceManager struct {
	db *gorm.DB
}

func NewDbServiceManager(db *gorm.DB) *DbServiceManager {
	return &DbServiceManager{
		db: db,
	}
}

func (sm *DbServiceManager) Find(serviceName string) ([]Location, error) {
	var service storage.Service
	if err := sm.db.Where("Name = ?", serviceName).First(&service).Error; err != nil {
		return nil, err
	}

	return []Location{Location{
		Address: service.Address,
		Port:    service.Port,

		Transport: service.Transport,
		Protocol:  service.Protocol,

		Weight:   service.Weight,
		Priority: service.Priority,
	},
	}, nil
}

func (sm *DbServiceManager) Register(serviceName string, location Location) error {
	var service storage.Service
	return sm.db.Where(&storage.Service{
		Name:      serviceName,
		Address:   location.Address,
		Transport: location.Transport,
		Protocol:  location.Protocol,
	}).Assign(&storage.Service{
		Port:     location.Port,
		Priority: location.Priority,
		Weight:   location.Weight,
	}).FirstOrCreate(&service).Error
}
