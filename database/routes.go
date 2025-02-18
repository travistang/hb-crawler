package database

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

const (
	RoutesDBMigration_13_21_23 = `
		ALTER TABLE routes ADD COLUMN crawledAt INTEGER DEFAULT NULL;
		ALTER TABLE routes ADD COLUMN elevation_loss REAL DEFAULT NULL;
		ALTER TABLE routes ADD COLUMN elevation_gain REAL DEFAULT NULL;
		ALTER TABLE routes ADD COLUMN t1_distance REAL DEFAULT NULL;
		ALTER TABLE routes ADD COLUMN t2_distance REAL DEFAULT NULL;
		ALTER TABLE routes ADD COLUMN t3_distance REAL DEFAULT NULL;
		ALTER TABLE routes ADD COLUMN t4_distance REAL DEFAULT NULL;
		ALTER TABLE routes ADD COLUMN t5_distance REAL DEFAULT NULL;
		ALTER TABLE routes ADD COLUMN t6_distance REAL DEFAULT NULL;
	`
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
	CrawledAt     *int
	ElevationGain *float32
	ElevationLoss *float32
	T1_Distance   *float32
	T2_Distance   *float32
	T3_Distance   *float32
	T4_Distance   *float32
	T5_Distance   *float32
	T6_Distance   *float32
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

	migrations := []string{RoutesDBMigration_13_21_23}
	for _, query := range migrations {
		// ignore error as a;
		repo.db.Exec(query)
	}
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
func (repo *RouteRepository) GetNextId() (*int, error) {
	query := `
		SELECT MAX(id) FROM routes
		WHERE crawledAt IS NOT NULL
	`

	routeId := 0
	if err := repo.Conn().QueryRow(query).Scan(&routeId); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	routeId += 1
	return &routeId, nil
}
