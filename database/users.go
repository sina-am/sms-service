package database

import (
	"main/entities"

	"gorm.io/gorm"
)

func GetUserById(id uint) (entities.User, error) {
	var user entities.User
	db := GetConnection()
	result := db.Where("id = ?", id).Find(&user)
	if result.RowsAffected != 1 {
		return user, gorm.ErrRecordNotFound
	}
	return user, nil
}

func AuthenticateUser(username, password string) (entities.User, error) {
	var user entities.User
	db := GetConnection()
	result := db.Where("username = ? and password = ?", username, password).Find(&user)
	if result.RowsAffected != 1 {
		return user, gorm.ErrRecordNotFound
	}
	return user, nil
}
