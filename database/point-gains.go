package database

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type PointGainsRepository struct {
	db *sql.DB
}

type PointGainRecord struct {
	EventId          int
	RoutePoints      int
	UserId           int
	UserPointsBefore int
	UserPointsAfter  *int
}

func CreatePointGainsRepository(db *sql.DB) *PointGainsRepository {
	return &PointGainsRepository{db: db}
}

func (repo *PointGainsRepository) Conn() *sql.DB {
	return repo.db
}

func (repo *PointGainsRepository) Migrate() error {
	log.Debugf("Migrating points gain repository...")
	query := `
		CREATE TABLE IF NOT EXISTS pointsGain(
			eventId INTEGER NOT NULL,
			userId INTEGER NOT NULL,
			routePoints INTEGER NOT NULL,
			pointsBefore INTEGER NOT NULL,
			pointsAfter INTEGER,

			PRIMARY KEY (eventId, userId)
		);
	`

	_, err := repo.db.Exec(query)
	return err
}

func (repo *PointGainsRepository) CreatePointsGainEntry(
	pointsGain *PointGainRecord,
) error {
	query := `
		INSERT OR IGNORE INTO pointsGain(
			eventId, userId, routePoints, pointsBefore, pointsAfter
		) VALUES(?, ?, ?, ?)
		ON CONFLICT(eventId, userId) DO UPDATE SET
			pointsBefore=excluded.pointsBefore
	`
	_, err := PrepareAndExecute(
		repo.Conn(), query,
		pointsGain.EventId, pointsGain.UserId, pointsGain.RoutePoints,
		pointsGain.UserPointsBefore, pointsGain.UserPointsAfter,
	)
	return err
}

func (repo *PointGainsRepository) UpdatePointsGainEntry(
	pointsGain *PointGainRecord,
) error {
	if pointsGain.UserPointsAfter == nil {
		return fmt.Errorf(
			"failed to update points gain entry without user points after:\n EventID: %d\nUserId: %d",
			pointsGain.EventId, pointsGain.UserId,
		)
	}

	query := `
		UPDATE pointsGain
		SET pointsAfter=?
		WHERE eventId=? AND userId=? AND pointsAfter IS NOT NULL
	`
	_, err := PrepareAndExecute(
		repo.Conn(), query,
		pointsGain.UserPointsAfter,
		pointsGain.EventId, pointsGain.UserId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PointGainsRepository) GetPointGainsByEventId(id int) (*[]PointGainRecord, error) {
	query := `
		SELECT eventId, userId, routePoints, pointsBefore, pointsAfter
		FROM pointsGain
		WHERE eventId=?
	`

	rows, err := repo.Conn().Query(query, id)
	if err != nil {
		return nil, err
	}

	records := []PointGainRecord{}
	for rows.Next() {
		var nextRecord PointGainRecord
		if err := rows.Scan(
			&nextRecord.EventId, &nextRecord.UserId,
			&nextRecord.RoutePoints, &nextRecord.UserPointsBefore, &nextRecord.UserPointsAfter,
		); err != nil {
			return nil, err
		}

		records = append(records, nextRecord)
	}

	return &records, nil
}
