package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func InitializeDatabase(connectionString string) (gorm.DB, error) {
	database, err := gorm.Open("postgres", connectionString)
	if err := database.DB().Ping(); err != nil {
		return gorm.DB{}, err
	}

	if err != nil {
		return gorm.DB{}, err
	}

	//database.LogMode(true)
	initializeTables(database)
	return database, nil
}

func initializeTables(database gorm.DB) {
	database.
		AutoMigrate(&Game{}).
		AutoMigrate(&Server{}).
		AutoMigrate(&Player{}).
		AutoMigrate(&Round{}).
		AutoMigrate(&WordList{}).
		AutoMigrate(&Setting{})
}
