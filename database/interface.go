package database

import (
	"main/entities"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	GetProviderById(int) (*entities.Provider, error)
	GetProviderByUsername(string) (*entities.Provider, error)
	FoundProviderByUsername(string) (bool, error)
	GetAllProviders() ([]*entities.Provider, error)
	CreateProvider(*entities.Provider) error
}
