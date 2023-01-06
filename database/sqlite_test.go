package database

import (
	"database/sql"
	"log"
	"main/entities"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func newFakeStorage() Storage {
	storage, err := NewSqliteStorage(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	return storage
}

func newFakeProvider() *entities.Provider {
	return &entities.Provider{
		Username:           "test",
		Password:           "test",
		PhoneNumber:        "test",
		Success:            10,
		Failed:             9,
		InvalidCredential:  false,
		InSufficientCredit: false,
	}
}

func newFakeProviders() []*entities.Provider {
	return []*entities.Provider{
		{
			Username:           "test",
			Password:           "test",
			PhoneNumber:        "test",
			Success:            1,
			Failed:             9,
			InvalidCredential:  true,
			InSufficientCredit: false,
		},
		{
			Username:           "test2",
			Password:           "test2",
			PhoneNumber:        "test2",
			Success:            10,
			Failed:             10,
			InvalidCredential:  false,
			InSufficientCredit: false,
		},
		{
			Username:           "test3",
			Password:           "test3",
			PhoneNumber:        "test3",
			Success:            7,
			Failed:             1,
			InvalidCredential:  false,
			InSufficientCredit: false,
		},
		{
			Username:           "test4",
			Password:           "test4",
			PhoneNumber:        "test4",
			Success:            150,
			Failed:             14,
			InvalidCredential:  false,
			InSufficientCredit: true,
		},
	}
}

func TestCreateProvider(t *testing.T) {
	p := newFakeProvider()
	storage := newFakeStorage()

	err := storage.CreateProvider(p)
	if err != nil {
		t.Error(err)
	}
}

func TestGetProviderById(t *testing.T) {
	p := newFakeProvider()
	storage := newFakeStorage()
	err := storage.CreateProvider(p)
	if err != nil {
		t.Error(err)
	}

	_, err = storage.GetProviderById(0)
	if err != sql.ErrNoRows {
		t.Errorf("expected not found error")
	}

	_, err = storage.GetProviderById(1)
	if err != nil {
		t.Error(err)
	}
}

func TestGetProviderByUsername(t *testing.T) {
	p := newFakeProvider()
	storage := newFakeStorage()
	err := storage.CreateProvider(p)
	if err != nil {
		t.Error(err)
	}

	_, err = storage.GetProviderByUsername("tee")
	if err != sql.ErrNoRows {
		t.Errorf("expected not found error")
	}

	p2, err := storage.GetProviderByUsername(p.Username)
	if err != nil {
		t.Error(err)
	}

	p.Id = p2.Id // ID's aren't set
	if !cmp.Equal(p, p2) {
		t.Errorf("%v, %v are not equal", p, p2)
	}
}

func TestGetAvailableProviders(t *testing.T) {
	storage := newFakeStorage()
	providers := newFakeProviders()
	for _, p := range providers {
		if err := storage.CreateProvider(p); err != nil {
			t.Error(err)
		}
	}

	sortedProviders, err := storage.GetAvailableProviders()
	if err != nil {
		t.Error(err)
	}
	if len(sortedProviders) != 2 {
		t.Errorf("wrong sort %v", sortedProviders)
	}
	if cmp.Equal(sortedProviders[0], providers[2]) {
		t.Errorf("wrong sort %v", sortedProviders)
	}
}

func TestUpdateProvider(t *testing.T) {
	storage := newFakeStorage()
	providers := newFakeProviders()
	for _, p := range providers {
		if err := storage.CreateProvider(p); err != nil {
			t.Error(err)
		}
	}

	if _, err := storage.GetProviderByUsername("updatedusername"); err != sql.ErrNoRows {
		t.Errorf("exists")
	}
	p := newFakeProvider()
	p.Id = 2
	err := storage.UpdateProvider(p)
	if err != nil {
		t.Error(err)
	}
	p2, err := storage.GetProviderById(p.Id)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(p, p2) {
		t.Error("not match")
	}
}
