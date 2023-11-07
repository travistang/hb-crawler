package database

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	PasswordHashEnvVariable = "HB_PASSWORD_HASH"
	DefaultPasswordHashKey  = "4b1532e0acb08de358c4e7a8619549426c864524093d242454e9948695fc438a"
)

type LoginCredentialRepository struct {
	db *sql.DB
}

type LoginCredentialRecord struct {
	SessionId, Username string
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateLoginCredentialRepository(db *sql.DB) *LoginCredentialRepository {
	return &LoginCredentialRepository{db: db}
}

func (repo *LoginCredentialRepository) Migrate() error {
	log.Debugf("Migrating login credentials repository...")
	query := `
		CREATE TABLE IF NOT EXISTS accounts(
			username TEXT PRIMARY KEY,
			password TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS credentials(
			username TEXT PRIMARY KEY,
			sessionid TEXT NOT NULL,
			date INTEGER DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := repo.db.Exec(query)
	return err
}

func (repo *LoginCredentialRepository) Conn() *sql.DB {
	return repo.db
}

func (repo *LoginCredentialRepository) GetCredential(record *LoginCredentialRecord, username string, oldestDate int64) error {
	query := `
		SELECT sessionid, username
		FROM credentials
		WHERE username=? AND date>=?;
	`
	log.Debugf("Getting credentials for username %s", username)
	if err := repo.Conn().
		QueryRow(query, username, oldestDate).
		Scan(&(record.SessionId), &(record.Username)); err != nil {
		if err == sql.ErrNoRows {
			log.Debugf("Credentials not found for user %s", username)
			return nil
		}
		return err
	}
	log.Debugf("Returning credential session %s", record.SessionId)

	return nil
}

func (repo *LoginCredentialRepository) SaveCredential(credential *LoginCredentialRecord) error {
	query := `
		INSERT OR IGNORE INTO credentials(sessionid, username, date)
		VALUES(?, ?, ?)
		ON CONFLICT(username) DO UPDATE SET 
			sessionid=excluded.sessionid,
			date=excluded.date;
	`
	_, err := PrepareAndExecute(
		repo.db, query,
		credential.SessionId, credential.Username, time.Now().Unix(),
	)
	return err
}

func (repo *LoginCredentialRepository) GetAvailableAccount() (*Account, error) {
	query := `
		SELECT username, password
		FROM accounts
		WHERE username IN (select username FROM accounts ORDER BY RANDOM() LIMIT 1)
	`
	var account Account
	if err := repo.Conn().QueryRow(query).Scan(&account.Username, &account.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func (repo *LoginCredentialRepository) GetAllAvailableAccounts() ([]Account, error) {
	query := `
		SELECT username, password
		FROM accounts
	`
	accounts := []Account{}
	rows, err := repo.Conn().Query(query)
	if err != nil {
		return accounts, err
	}
	var account Account
	for rows.Next() {
		if err := rows.Scan(&account.Username, &account.Password); err != nil {
			return []Account{}, err
		}
		accounts = append(accounts, Account{
			Username: account.Username,
		})
	}

	return accounts, nil
}

func (repo *LoginCredentialRepository) CreateAccount(account *Account) error {
	query := `
		INSERT INTO accounts(username, password)
		VALUES (?, ?)
	`
	if _, err := PrepareAndExecute(
		repo.Conn(), query,
		account.Username, account.Password,
	); err != nil {
		return err
	}

	return nil
}
