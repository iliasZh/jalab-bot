package todays_jalabs

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"jalabs.kz/bot/model/db"
)

const (
	queryCreate = `
	INSERT INTO todays_jalabs (group_chat_id, user_id, created_at)
	VALUES (:group_chat_id, :user_id, NOW() AT TIME ZONE 'utc')
	RETURNING id, group_chat_id, user_id, created_at;`
)

type Repo struct {
	db         *sqlx.DB
	stmtCreate *sqlx.NamedStmt
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

	return Repo{
		db:         db,
		stmtCreate: stmtCreate,
	}, nil
}

func (r Repo) Create(ctx context.Context, tj db.TodaysJalab) (created db.TodaysJalab, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("creating todays jalab %+v: %w", tj, err)
	}()

	errGet := r.stmtCreate.GetContext(ctx, &created, tj)
	if errGet != nil {
		return db.TodaysJalab{}, fmt.Errorf("executing query: %w", errGet)
	}

	return created, nil
}
