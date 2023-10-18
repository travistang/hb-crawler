package database

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type UserRepository struct {
	db *sql.DB
}

type UserRecord struct {
	Id             int
	Name, LastName string
}

func CreateUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Conn() *sql.DB {
	return repo.db
}

func (repo *UserRepository) Migrate() error {
	log.Debugf("Migrating user repository...")
	query := `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			lastname INTEGER NOT NULL
		);
	`
	_, err := repo.db.Exec(query)
	return err
}

func (repo *UserRepository) GetUser(id int) (*UserRecord, error) {
	query := `
		SELECT 
			id, name, lastname 
		FROM users
		WHERE id=?
	`
	rows, err := PrepareAndQuery(repo.Conn(), query, id)
	if err != nil {
		return nil, err
	}

	var user UserRecord
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.LastName)
		break
	}
	return &user, nil
}
