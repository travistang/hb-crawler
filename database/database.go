package database

import (
	"database/sql"
)

func InitializeDatabase(pathName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", pathName)
	if err != nil {
		return nil, err
	}

	repositories := []Repository{
		&EventRepository{db: db},
		&UserRepository{db: db},
	}

	for _, repo := range repositories {
		if err := repo.Migrate(); err != nil {
			db.Close()
			return nil, err
		}
	}

	return db, nil
}
