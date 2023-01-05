package database

import (
	"database/sql"
	"main/entities"
)

type sqliteStorage struct {
	db *sql.DB
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS providers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			type VARCHAR(255) NOT NULL,
			phone_number VARCHAR(255) NOT NULL
		);
	`)
	return err
}

func NewSqliteStorage(connStr string) (*sqliteStorage, error) {
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		return nil, err
	}
	return &sqliteStorage{db: db}, nil
}

func (s *sqliteStorage) GetProviderById(id int) (*entities.Provider, error) {
	provider := &entities.Provider{}
	row := s.db.QueryRow(`SELECT id, username, password, phone_number, type FROM providers WHERE id = $1`, id)
	err := row.Scan(
		&provider.Id,
		&provider.Username,
		&provider.Password,
		&provider.PhoneNumber,
		&provider.Type,
	)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *sqliteStorage) GetProviderByUsername(username string) (*entities.Provider, error) {
	provider := &entities.Provider{}
	row := s.db.QueryRow(
		`SELECT id, username, password, phone_number, type FROM providers WHERE username = $1`, username)
	err := row.Scan(
		&provider.Id,
		&provider.Username,
		&provider.Password,
		&provider.PhoneNumber,
		&provider.Type,
	)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

/*
Checks if the given provider exists or not

	returns true if exist
*/
func (s *sqliteStorage) FoundProviderByUsername(username string) (bool, error) {
	err := s.db.QueryRow(
		`SELECT username FROM providers WHERE username = $1`, username).Scan(&username)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *sqliteStorage) GetAllProviders() ([]*entities.Provider, error) {
	providers := []*entities.Provider{}
	rows, err := s.db.Query(`SELECT id, username, password, phone_number, type FROM providers`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		provider := &entities.Provider{}
		err := rows.Scan(
			&provider.Id,
			&provider.Username,
			&provider.Password,
			&provider.PhoneNumber,
			&provider.Type,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func (s sqliteStorage) CreateProvider(p *entities.Provider) error {
	_, err := s.db.Exec(
		`INSERT INTO providers(username, password, phone_number, type) VALUES ($1,$2,$3,$4)`,
		p.Username, p.Password, p.PhoneNumber, p.Type,
	)
	return err
}
