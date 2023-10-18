package todays_jalabs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"jalabs.kz/bot/model/db"
	"jalabs.kz/bot/storage/postgres/error_code"
)

var (
	ErrNotFound      = fmt.Errorf("today's jalab not found")
	ErrAlreadyExists = fmt.Errorf("today's jalab already exists")
)

const (
	queryCreate = `
	INSERT INTO todays_jalabs (group_chat_id, user_id, created_at)
	VALUES (:group_chat_id, :user_id, CURRENT_DATE)
	RETURNING id, group_chat_id, user_id, created_at;`

	queryGet = `
	SELECT id, group_chat_id, user_id, created_at
	FROM todays_jalabs
	WHERE group_chat_id = :group_chat_id AND created_at = CURRENT_DATE;`
)

type Repo struct {
	db         *sqlx.DB
	stmtCreate *sqlx.NamedStmt
	stmtGet    *sqlx.NamedStmt
}

func New(ctx context.Context, db *sqlx.DB) (r Repo, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("new todays jalabs repo: %w", err)
	}()

	stmtCreate, errPrepareCreate := db.PrepareNamedContext(ctx, queryCreate)
	if errPrepareCreate != nil {
		return Repo{}, fmt.Errorf("preparing create stmt: %w", errPrepareCreate)
	}

	stmtGet, errPrepareGet := db.PrepareNamedContext(ctx, queryGet)
	if errPrepareGet != nil {
		return Repo{}, fmt.Errorf("preparing get stmt: %w", errPrepareGet)
	}

	return Repo{
		db:         db,
		stmtCreate: stmtCreate,
		stmtGet:    stmtGet,
	}, nil
}

func (r Repo) Create(ctx context.Context, tj db.TodaysJalab) (created db.TodaysJalab, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("creating todays jalab %+v: %w", tj, err)
	}()

	errCreate := r.stmtCreate.GetContext(ctx, &created, tj)
	if error_code.GetFrom(errCreate) == error_code.UniqueConstraintViolation {
		return db.TodaysJalab{}, ErrAlreadyExists
	}
	if errCreate != nil {
		return db.TodaysJalab{}, fmt.Errorf("executing query: %w", errCreate)
	}

	return created, nil
}

func (r Repo) Get(ctx context.Context, tj db.TodaysJalab) (gotten db.TodaysJalab, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("getting todays jalab of group with id %d: %w", tj.GroupChatID, err)
	}()

	errGet := r.stmtGet.GetContext(ctx, &gotten, tj)
	if errors.Is(errGet, sql.ErrNoRows) {
		return db.TodaysJalab{}, ErrNotFound
	}
	if errGet != nil {
		return db.TodaysJalab{}, fmt.Errorf("executing query: %w", errGet)
	}

	return gotten, nil
}

func (r Repo) CreateOrGet(
	ctx context.Context,
	tj db.TodaysJalab,
) (retrieved db.TodaysJalab, err error) {
	jalab, errCreate := r.Create(ctx, tj)
	switch {
	case errCreate == nil:
		return jalab, nil
	case errors.Is(errCreate, ErrAlreadyExists):
		return r.Get(ctx, tj)
	default:
		return db.TodaysJalab{}, errCreate
	}
}
