package database

import (
	"database/sql"
)

type EventRepository struct {
	db *sql.DB
}

type EventRecord struct {
	Id, RouteId, OrganizerId, Points int
	Title                            string
}

func (repo *EventRepository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS events(
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			routeId INTEGER NOT NULL,
			points INTEGER,
			organizerId INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS eventParticipations(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			eventId INTEGER,
			userId INTEGER,
			FOREIGN KEY(eventId) REFERENCES events(id),
			FOREIGN KEY(userId) REFERENCES users(id)
		)
	`
	_, err := repo.db.Exec(query)
	return err
}

func (repo *EventRepository) Conn() *sql.DB {
	return repo.db
}

// func (repo *EventRepository) GetEventById(id int) (*EventRecord, error) {
// 	query := `
// 		SELECT
// 			id, title, points,
// 			users.id
// 		FROM
// 			events
// 			LEFT JOIN eventParticipations ON events.id = eventParticipations.eventId
// 			LEFT JOIN users ON users.id = eventParticipations.userId
// 		WHERE id=?
// 	`
// 	rows, err := PrepareAndExecute(repo.Conn(), query, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	event := EventRecord{}

// 	for rows.Next() {
// 		rows.Scan(&event.Id)
// 	}

// }
