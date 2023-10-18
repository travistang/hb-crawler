package database

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

type EventRepository struct {
	db *sql.DB
}

type EventRecord struct {
	Id, RouteId, OrganizerId, Points int
	Date                             time.Time
	Title                            string
}

func (repo *EventRepository) Migrate() error {
	log.Debugf("Migrating events repository...")

	query := `
		CREATE TABLE IF NOT EXISTS events(
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			routeId INTEGER NOT NULL,
			date INTEGER DEFAULT CURRENT_TIMESTAMP,
			organizerId INTEGER NOT NULL
		);
	`
	_, err := repo.db.Exec(query)
	return err
}

func (repo *EventRepository) Conn() *sql.DB {
	return repo.db
}

func CreateEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}
