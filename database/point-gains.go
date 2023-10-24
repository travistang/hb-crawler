package database

import (
	"database/sql"
	"fmt"
	"time"

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
	EventDate        int64
}

type ReducedPointGainRecord struct {
	RoutePoints      int `json:"route_points"`
	UserPointsBefore int `json:"points_before"`
	UserPointsAfter  int `json:"points_after"`
}

type PointGainsQuery struct {
	Limit int
	Skip  int
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
			eventDate INTEGER NOT NULL,

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
			eventId, userId, routePoints, pointsBefore, pointsAfter, eventDate
		) VALUES(?, ?, ?, ?, ?, ?)
		ON CONFLICT(eventId, userId) DO UPDATE SET
			pointsBefore=excluded.pointsBefore
	`
	_, err := PrepareAndExecute(
		repo.Conn(), query,
		pointsGain.EventId, pointsGain.UserId, pointsGain.RoutePoints,
		pointsGain.UserPointsBefore, pointsGain.UserPointsAfter,
		pointsGain.EventDate,
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
		WHERE eventId=? AND userId=? AND pointsAfter IS NULL
	`
	_, err := PrepareAndExecute(
		repo.Conn(), query,
		*pointsGain.UserPointsAfter,
		pointsGain.EventId, pointsGain.UserId,
	)

	if err != nil {
		return err
	}

	return nil
}

func extractRowToRecords(rows *sql.Rows) (*[]PointGainRecord, error) {
	records := []PointGainRecord{}
	for rows.Next() {
		var nextRecord PointGainRecord
		if err := rows.Scan(
			&nextRecord.EventId, &nextRecord.UserId,
			&nextRecord.RoutePoints, &nextRecord.UserPointsBefore, &nextRecord.UserPointsAfter,
			&nextRecord.EventDate,
		); err != nil {
			return nil, err
		}

		records = append(records, nextRecord)
	}

	return &records, nil
}

func (repo *PointGainsRepository) GetPointGainsByEventId(id int) (*[]PointGainRecord, error) {
	query := `
		SELECT 
			eventId, userId, routePoints, pointsBefore, pointsAfter, eventDate
		FROM pointsGain
		WHERE eventId=?
	`

	rows, err := repo.Conn().Query(query, id)
	if err != nil {
		return nil, err
	}

	return extractRowToRecords(rows)
}

func (repo *PointGainsRepository) GetAllPointGains(queryParams *PointGainsQuery) (*[]PointGainRecord, error) {
	params := PointGainsQuery{
		Skip:  0,
		Limit: -1,
	}
	if queryParams != nil {
		params.Limit = queryParams.Limit
		params.Skip = queryParams.Skip
	}

	query := `
		SELECT 
			eventId, userId, routePoints, pointsBefore, pointsAfter, eventDate
		FROM pointsGain
		ORDER BY eventDate ASC
		OFFSET ?
		LIMIT ?
	`
	rows, err := repo.Conn().Query(query, params.Skip, params.Limit)
	if err != nil {
		return nil, err
	}

	return extractRowToRecords(rows)
}

func (repo *PointGainsRepository) GetDanglingPointsGainEntryToday(targetHour time.Time) (*[]PointGainRecord, error) {
	query := `
		SELECT
			eventId, userId, routePoints, pointsBefore, pointsAfter, eventDate
		FROM pointsGain
		WHERE 
			pointsAfter IS NULL AND
			eventDate >= ? AND eventDate <= ?
	`
	dayStart := time.Date(
		targetHour.Year(), targetHour.Month(), targetHour.Day(), 0,
		0, 0, 0, targetHour.Location(),
	)
	hourStart := time.Date(
		targetHour.Year(), targetHour.Month(), targetHour.Day(), targetHour.Hour(),
		0, 0, 0, targetHour.Location(),
	)
	rows, err := repo.Conn().Query(query, dayStart.Unix(), hourStart.Unix())
	if err != nil {
		return nil, err
	}

	return extractRowToRecords(rows)
}

func (repo *PointGainsRepository) GetValidPointsGainEntry(desiredLimit *int) (*[]ReducedPointGainRecord, error) {
	query := `
		SELECT
			routePoints, pointsBefore, pointsAfter
		FROM pointsGain
		WHERE
			pointsAfter IS NOT NULL AND
			pointsBefore < pointsAfter
		ORDER BY eventDate DESC
		LIMIT ?
	`
	limit := -1
	if desiredLimit != nil {
		limit = *desiredLimit
	}
	rows, err := repo.Conn().Query(query, limit)
	if err != nil {
		return nil, err
	}

	records := []ReducedPointGainRecord{}
	for rows.Next() {
		var nextRecord ReducedPointGainRecord
		if err := rows.Scan(
			&nextRecord.RoutePoints, &nextRecord.UserPointsBefore, &nextRecord.UserPointsAfter,
		); err != nil {
			return nil, err
		}

		records = append(records, nextRecord)
	}

	return &records, nil
}
