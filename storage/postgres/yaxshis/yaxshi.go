package yaxshis

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"jalabs.kz/bot/model/db"
	"jalabs.kz/bot/storage/postgres/error_code"
)

var ErrNoSuchJalab = fmt.Errorf("no such jalab")

const (
	queryCreate = `
	INSERT INTO yaxshis (group_chat_id, user_id, count, created_at, updated_at)
	VALUES (:group_chat_id, :user_id, 1, NOW() AT TIME ZONE 'utc', NOW() AT TIME ZONE 'utc')
	ON CONFLICT (group_chat_id, user_id)
	    DO UPDATE SET count = yaxshis.count + 1, updated_at = NOW() AT TIME ZONE 'utc'
	RETURNING id, group_chat_id, user_id, count, created_at, updated_at;`
)

type Repo struct {
	db *sqlx.DB

	stmtCreate *sqlx.NamedStmt
}

func New(ctx context.Context, db *sqlx.DB) (r Repo, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("new yaxshis repo: %w", err)
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

func (r Repo) Create(ctx context.Context, y db.Yaxshi) (created db.Yaxshi, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("creating yaxshi %+v: %w", y, err)
	}()

	errCreate := r.stmtCreate.GetContext(ctx, &created, y)
	if error_code.GetFrom(errCreate) == error_code.ForeignKeyViolation {
		return db.Yaxshi{}, ErrNoSuchJalab
	}
	if errCreate != nil {
		return db.Yaxshi{}, fmt.Errorf("executing query: %w", errCreate)
	}

	return created, nil
}
