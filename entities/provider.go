package entities

import "gorm.io/gorm"

type ProviderType string

const (
	Melipayamak ProviderType = "Melipayamak"
	Celery      ProviderType = "Celery"
)

type Provider struct {
	gorm.Model
	ID          uint         `json:"id"`
	Type        ProviderType `json:"type" binding:"required"`
	Username    string       `json:"username" binding:"required"`
	Password    string       `json:"password" binding:"required"`
	PhoneNumber string       `json:"phone_number" binding:"required"`
}
