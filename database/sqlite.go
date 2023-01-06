package database

import (
	"database/sql"
	"main/entities"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteStorage struct {
	db *sql.DB
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS providers (
			id 					INTEGER PRIMARY KEY AUTOINCREMENT,
			username			VARCHAR(255) NOT NULL,
			password			VARCHAR(255) NOT NULL,
			type 				VARCHAR(255) NOT NULL,
			phone_number 		VARCHAR(255) NOT NULL,
			success 			INTEGER DEFAULT 0,
			failed 				INTEGER DEFAULT 0,
			invalid_credential 	BOOLEAN DEFAULT false,
			insufficient_credit BOOLEAN DEFAULT false
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

type Scanner interface {
	Scan(dest ...any) error
}

func (s *sqliteStorage) scanIntoProvider(row Scanner) (*entities.Provider, error) {
	provider := &entities.Provider{}
	err := row.Scan(
		&provider.Id,
		&provider.Username,
		&provider.Password,
		&provider.PhoneNumber,
		&provider.Type,
		&provider.Success,
		&provider.Failed,
		&provider.InvalidCredential,
		&provider.InSufficientCredit,
	)
	return provider, err
}

func (s *sqliteStorage) scanIntoProviders(rows *sql.Rows) ([]*entities.Provider, error) {
	providers := []*entities.Provider{}
	for rows.Next() {
		provider, err := s.scanIntoProvider(rows)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func (s *sqliteStorage) GetProviderById(id int) (*entities.Provider, error) {
	row := s.db.QueryRow(`
		SELECT
			id, username, password,
			phone_number, type,
			success, failed, invalid_credential,
			insufficient_credit
			FROM providers
			WHERE id = $1`,
		id,
	)
	return s.scanIntoProvider(row)
}

func (s *sqliteStorage) GetProviderByUsername(username string) (*entities.Provider, error) {
	row := s.db.QueryRow(`
		SELECT
			id, username, password,
			phone_number, type,
			success, failed, invalid_credential,
			insufficient_credit
		FROM providers
		WHERE username = $1`,
		username,
	)
	return s.scanIntoProvider(row)
}

/*
Checks if the given provider exists or not

	returns true if exist
*/
func (s *sqliteStorage) ExistProviderByUsername(username string) (bool, error) {
	err := s.db.QueryRow(`
		SELECT 
			username 
		FROM providers 
		WHERE username = $1`,
		username).Scan(&username)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *sqliteStorage) GetAllProviders() ([]*entities.Provider, error) {
	rows, err := s.db.Query(`
		SELECT 
			id, username, password, 
			phone_number, type,
			success, failed, invalid_credential,
			insufficient_credit 
		FROM providers`,
	)
	if err != nil {
		return nil, err
	}
	return s.scanIntoProviders(rows)
}

/* Sort providers base on availability */
func (s *sqliteStorage) GetAvailableProviders() ([]*entities.Provider, error) {
	rows, err := s.db.Query(`
		SELECT 
			id, username, password, 
			phone_number, type,
			success, failed, invalid_credential,
			insufficient_credit 
		FROM providers
		WHERE 
			invalid_credential = false AND 
			insufficient_credit = false
		ORDER BY failed ASC, success DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	return s.scanIntoProviders(rows)
}
func (s sqliteStorage) CreateProvider(p *entities.Provider) error {
	_, err := s.db.Exec(`
		INSERT INTO 
			providers(username, password, phone_number, type, success, 
				failed, invalid_credential, insufficient_credit) 
		VALUES ($1,$2,$3,$4, $5, $6, $7, $8)`,
		p.Username, p.Password, p.PhoneNumber,
		p.Type, p.Success, p.Failed,
		p.InvalidCredential, p.InSufficientCredit,
	)
	return err
}

func (s sqliteStorage) UpdateProvider(p *entities.Provider) error {
	_, err := s.db.Exec(`
		UPDATE providers
		SET 
			username = $1,
			password = $2,
			phone_number = $3,
			type = $4,
			success = $6, 
			failed = $7, 
			invalid_credential = $8,
			insufficient_credit = $9 
		WHERE id = $10
	`,
		p.Username, p.Password, p.PhoneNumber, p.Type,
		p.Success, p.Failed, p.InvalidCredential,
		p.InSufficientCredit, p.Id,
	)
	return err
}
