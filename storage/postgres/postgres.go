package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	dsn = "host=localhost port=5432 sslmode=disable user=admin password=admin dbname=jalab_bot"

	SQLStateForeignKeyViolation       = "23503"
	SQLStateUniqueConstraintViolation = "23505"
)

func New() (*sqlx.DB, error) {
	db, errConnect := sqlx.Connect("postgres", dsn)
	if errConnect != nil {
		return nil, fmt.Errorf("connecting to database: %w", errConnect)
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)
	db.SetConnMaxIdleTime(60)

	return db, nil
}
