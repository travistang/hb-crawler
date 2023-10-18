package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func InitializeDatabase(pathName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", pathName)

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
