package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"jalabs.kz/bot/model/db"
)

var ErrNotFound = fmt.Errorf("user not found")

const (
	queryCreate = `
	INSERT INTO users (id, username, first_name, created_at, updated_at)
	VALUES (:id, :username, :first_name, NOW() AT TIME ZONE 'utc', NOW() AT TIME ZONE 'utc')
	ON CONFLICT (id) DO
	    UPDATE SET username = :username,
	               first_name = :first_name,
	               updated_at = NOW() AT TIME ZONE 'utc'
	RETURNING id, username, first_name, disable_mentions, created_at, updated_at;`

	queryToggleMentions = `
	UPDATE users
	SET disable_mentions = NOT disable_mentions
	WHERE id = :id
	RETURNING id, username, first_name, disable_mentions, created_at, updated_at;`

	queryGet = `
	SELECT id, username, first_name, disable_mentions, created_at, updated_at
	FROM users
	WHERE id = :id;`
)

type Repo struct {
	db *sqlx.DB

	stmtCreate         *sqlx.NamedStmt
	stmtGet            *sqlx.NamedStmt
	stmtToggleMentions *sqlx.NamedStmt
}

func New(ctx context.Context, db *sqlx.DB) (r Repo, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("new users repo: %w", err)
	}()

	stmtCreate, errPrepareCreate := db.PrepareNamedContext(ctx, queryCreate)
	if errPrepareCreate != nil {
		return Repo{}, fmt.Errorf("preparing create stmt: %w", errPrepareCreate)
	}

	stmtGet, errPrepareGet := db.PrepareNamedContext(ctx, queryGet)
	if errPrepareGet != nil {
		return Repo{}, fmt.Errorf("preparing get stmt: %w", errPrepareGet)
	}

	stmtToggleMentions, errPrepareToggleMentions := db.PrepareNamedContext(ctx, queryToggleMentions)
	if errPrepareToggleMentions != nil {
		return Repo{}, fmt.Errorf(
			"preparing toggle mention stmt: %w", errPrepareToggleMentions,
		)
	}

	return Repo{
		db:                 db,
		stmtCreate:         stmtCreate,
		stmtGet:            stmtGet,
		stmtToggleMentions: stmtToggleMentions,
	}, nil
}

func (r Repo) Create(ctx context.Context, u db.User) (created db.User, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("creating user %+v: %w", u, err)
	}()

	errCreate := r.stmtCreate.GetContext(ctx, &created, u)
	if errCreate != nil {
		return db.User{}, fmt.Errorf("executing query: %w", errCreate)
	}

	return created, nil
}

func (r Repo) Get(ctx context.Context, u db.User) (gotten db.User, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("getting user %+v: %w", u, err)
	}()

	errGet := r.stmtGet.GetContext(ctx, &gotten, u)
	if errGet != nil {
		return db.User{}, fmt.Errorf("executing query: %w", errGet)
	}

	return gotten, nil
}

func (r Repo) ToggleMentions(ctx context.Context, u db.User) (updated db.User, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("toggling mentions for user %+v: %w", u, err)
	}()

	errToggle := r.stmtToggleMentions.GetContext(ctx, &updated, u)
	if errors.Is(errToggle, sql.ErrNoRows) {
		return db.User{}, ErrNotFound
	}
	if errToggle != nil {
		return db.User{}, errToggle
	}

	return updated, nil
}
