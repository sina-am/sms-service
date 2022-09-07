package database

import (
	"log"
	"main/entities"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db := GetConnection()
	db.AutoMigrate(&entities.Provider{})
	return db, nil
}

func GetConnection() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	return db
}
