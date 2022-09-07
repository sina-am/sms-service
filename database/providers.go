package database

import (
	"main/entities"
)

func GetProviderById(id uint) entities.Provider {
	var provider entities.Provider
	db := GetConnection()
	db.Where("id = ?", id).First(&provider)
	return provider
}

func GetProvidersByUserId(userId uint) []entities.Provider {
	var providers []entities.Provider
	db := GetConnection()
	db.Where("userId = ?", userId).Find(&providers)
	return providers
}

func CreateProvider(provider *entities.Provider) {
	db := GetConnection()
	db.Create(&provider)
}
