package database

import (
	"database/sql"
	hiking_buddies "hb-crawler/rating-gain/hiking-buddies"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

type UserRecord struct {
	Id             int
	Name, LastName string
}

type UserPointsHistory struct {
	Id, userId, points int
	date               time.Time
}

func CreateUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			lastname INTEGER NOT NULL,
		);

		CREATE TABLE IF NOT EXISTS userPointsHistory(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL,
			points INTEGER NOT NULL,
			date TEXT DEFAULT datetime('now'),

			FOREIGN KEY(userId) REFERENCES users(id)
		)
	`
	_, err := repo.db.Exec(query)
	return err
}

func (repo *UserRepository) Conn() *sql.DB {
	return repo.db
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

func (repo *UserRepository) GetPointsHistory(id int) ([]UserPointsHistory, error) {
	query := `
		SELECT 
			id, points, date
		FROM userPointsHistory
		WHERE userId=?
		ORDER BY date DESC
	`
	rows, err := PrepareAndQuery(repo.Conn(), query, id)
	if err != nil {
		return nil, err
	}

	histories := []UserPointsHistory{}

	for rows.Next() {
		history := UserPointsHistory{}

		err := rows.Scan(&history.Id, &history.points, &history.date)
		if err != nil {
			return nil, err
		}

		history.userId = id
		histories = append(histories, history)
	}

	return histories, nil
}

func (repo *UserRepository) CreateUser(user *hiking_buddies.User) (*int, error) {
	query := `
		INSERT OR IGNORE INTO users(
			id, name, lastname
		) VALUES(?, ?, ?)
	`
	newId, err := PrepareAndExecute(
		repo.Conn(), query,
		user.ID, user.Name, user.LastName,
	)

	if err != nil {
		return nil, err
	}

	id := int(*newId)
	return &id, nil
}
