package database

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

type LoginCredentialRepository struct {
	db *sql.DB
}

type LoginCredentialRecord struct {
	SessionId, Username string
}

func CreateLoginCredentialRepository(db *sql.DB) *LoginCredentialRepository {
	return &LoginCredentialRepository{db: db}
}

func (repo *LoginCredentialRepository) Migrate() error {
	log.Debugf("Migrating login credentials repository...")
	query := `
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
