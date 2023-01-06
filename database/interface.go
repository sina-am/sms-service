package database

import (
	"main/entities"
)

type Storage interface {
	GetProviderById(int) (*entities.Provider, error)
	GetProviderByUsername(string) (*entities.Provider, error)
	ExistProviderByUsername(string) (bool, error)
	GetAllProviders() ([]*entities.Provider, error)
	GetAvailableProviders() ([]*entities.Provider, error)

	CreateProvider(*entities.Provider) error
	UpdateProvider(*entities.Provider) error
}
