package cmd

import (
	"github.com/jmoiron/sqlx"
)

func DBConnect(cfg Config) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	// Let's connect to the DB.
	db, err = sqlx.Connect("postgres", cfg.PostgresURL)
	if err != nil {
		return db, err
	}

	// Make sure we've got a table to write to.
	db.MustExec(DBTable())

	// TODO: Setup partitions?

	return db, nil
}

// DBTable is the SQL that creates the main table in Postgres.
func DBTable() string {
	sql := `CREATE TABLE IF NOT EXISTS metrics (
			id serial PRIMARY KEY,
			captured_at timestamp,
			address text,
			response_time integer,
			status_code integer,
			regexp text,
			regexp_status text
		);`
	return sql
}
