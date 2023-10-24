package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const (
	dbFilename            = "db.sqlite"
	sameDirectoryPath     = "."
	dataDirectoryPath     = ".data"
	rootDataDirectoryPath = "/data"
	defaultDirectoryPath  = dataDirectoryPath + "/" + dbFilename
)

func locateDatabase() (*string, error) {
	paths := []string{
		sameDirectoryPath, dataDirectoryPath, rootDataDirectoryPath,
	}
	currentPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	currentPath = filepath.Dir(currentPath)
	for _, path := range paths {
		filename, err := url.JoinPath(currentPath, path, dbFilename)
		log.Infof("Searching database in location %s...", filename)
		if err != nil {
			continue
		}
		_, err = os.Stat(filename)
		if err == nil {
			return &filename, nil
		}
	}
	return nil, fmt.Errorf("no database file found")
}

func InitializeDatabase() (*sql.DB, error) {
	pathName, err := locateDatabase()
	if err != nil {
		return nil, err
	}
	log.Infof("Using database at location %s", *pathName)
	db, err := sql.Open("sqlite3", *pathName)

	if err != nil {
		return nil, err
	}

	repositories := []Repository{
		&EventRepository{db: db},
		&UserRepository{db: db},
		&LoginCredentialRepository{db: db},
		&RouteRepository{db: db},
		&PointGainsRepository{db: db},
	}

	for _, repo := range repositories {
		err := repo.Migrate()
		if err != nil {
			log.Errorf("Failed to migrate: %+v\n", err)

			db.Close()
			return nil, err
		}
	}

	return db, nil
}

type DatabaseRepository struct {
	User       *UserRepository
	Route      *RouteRepository
	Login      *LoginCredentialRepository
	Event      *EventRepository
	PointGains *PointGainsRepository
}

func GetRepository(db *sql.DB) *DatabaseRepository {
	user := CreateUserRepository(db)
	route := CreateRouteRepository(db)
	login := CreateLoginCredentialRepository(db)
	event := CreateEventRepository(db)
	pointGains := CreatePointGainsRepository(db)

	return &DatabaseRepository{
		User:       user,
		Route:      route,
		Login:      login,
		Event:      event,
		PointGains: pointGains,
	}
}
