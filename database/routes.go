package database

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type RouteRepository struct {
	db *sql.DB
}

type RouteRecord struct {
	Id, Elevation int
	Points        *int
	Distance      float32
	Name          string
	Scale         string
}

func (repo *RouteRepository) Migrate() error {
	log.Debugf("Migrating routes repository...")

	query := `
		CREATE TABLE IF NOT EXISTS routes(
			id INTEGER PRIMARY KEY,
			points INTEGER,
			elevation INTEGER NOT NULL,
			name TEXT NOT NULL,
			scale TEXT NOT NULL,
			distance REAL NOT NULL
		);
	`
	_, err := repo.db.Exec(query)
	return err
}

func (repo *RouteRepository) Conn() *sql.DB {
	return repo.db
}

func CreateRouteRepository(db *sql.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

func (repo *RouteRepository) SaveRoute(route *RouteRecord) (bool, error) {
	query := `
		INSERT OR IGNORE INTO routes(
			id, points, elevation, name, scale, distance
		) 
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			points=excluded.points;
	`
	db := repo.Conn()
	id, err := PrepareAndExecute(
		db, query,
		route.Id, route.Points, route.Elevation, route.Name, route.Scale, route.Distance,
	)
	if err != nil {
		return false, err
	}

	saved := *id > 0
	return saved, nil
}

func (repo *RouteRepository) GetRouteById(id int, record *RouteRecord) error {
	query := `
		SELECT id, name, points, elevation, scale, distance FROM routes
		WHERE id=?
	`
	var route RouteRecord
	if err := repo.Conn().QueryRow(query, id).Scan(
		&route.Id, &route.Name, route.Points, &route.Elevation, &route.Scale, &route.Distance,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	*record = route
	return nil
}
