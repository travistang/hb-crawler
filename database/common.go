package database

import "database/sql"

type Repository interface {
	Migrate() error
	Conn() *sql.DB
}

func PrepareAndQuery(conn *sql.DB, query string, args ...any) (*sql.Rows, error) {
	statement, err := conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(args)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func PrepareAndExecute(conn *sql.DB, query string, args ...any) (*int64, error) {
	statement, err := conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.Exec(args...)
	if err != nil {
		return nil, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &lastId, nil
}
