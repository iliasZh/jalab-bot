package db

import (
	"database/sql"
	"time"
)

type User struct {
	ID              int64          `db:"id"`
	Username        sql.NullString `db:"username"`
	FirstName       string         `db:"first_name"`
	DisableMentions bool           `db:"disable_mentions"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}

func (u User) UsernameString() string {
	if u.Username.Valid {
		if u.DisableMentions {
			return u.Username.String
		}

		return "@" + u.Username.String
	}

	return "@анонимус"
}
