package jalabs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"jalabs.kz/bot/model/db"
)

var (
	ErrAlreadyExists = fmt.Errorf("jalab already exists")
	ErrNotFound      = fmt.Errorf("jalab not found")
)

const (
	queryCreate = `
	INSERT INTO jalabs (group_chat_id, user_id, created_at, updated_at)
	VALUES (
	        :group_chat_id, :user_id,
	        NOW() AT TIME ZONE 'utc', NOW() AT TIME ZONE 'utc'
	)
	ON CONFLICT DO NOTHING
	RETURNING id, group_chat_id, user_id, created_at, updated_at;`

	queryGet = `
	SELECT id, group_chat_id, user_id, created_at, updated_at
	FROM jalabs
	WHERE group_chat_id = :group_chat_id AND user_id = :user_id;`

	queryGetAll = `
	SELECT id, group_chat_id, user_id, created_at, updated_at
	FROM jalabs
	WHERE group_chat_id = :group_chat_id;`
)

type Repo struct {
	db         *sqlx.DB
	stmtCreate *sqlx.NamedStmt
	stmtGet    *sqlx.NamedStmt
	stmtGetAll *sqlx.NamedStmt
}

func New(ctx context.Context, db *sqlx.DB) (r Repo, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("new jalabs repo: %w", err)
	}()

	stmtCreate, errPrepareCreate := db.PrepareNamedContext(ctx, queryCreate)
	if errPrepareCreate != nil {
		return Repo{}, fmt.Errorf("preparing create stmt: %w", errPrepareCreate)
	}

	stmtGet, errPrepareGet := db.PrepareNamedContext(ctx, queryGet)
	if errPrepareGet != nil {
		return Repo{}, fmt.Errorf("preparing get stmt: %w", errPrepareGet)
	}

	stmtGetAll, errPrepareGetAll := db.PrepareNamedContext(ctx, queryGetAll)
	if errPrepareGetAll != nil {
		return Repo{}, fmt.Errorf("preparing get all stmt: %w", errPrepareGetAll)
	}

	return Repo{
		db:         db,
		stmtCreate: stmtCreate,
		stmtGet:    stmtGet,
		stmtGetAll: stmtGetAll,
	}, nil
}

func (r Repo) Create(ctx context.Context, j db.Jalab) (created db.Jalab, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("creating jalab %+v: %w", j, err)
	}()

	errCreate := r.stmtCreate.GetContext(ctx, &created, j)
	if errors.Is(errCreate, sql.ErrNoRows) {
		return db.Jalab{}, ErrAlreadyExists
	}
	if errCreate != nil {
		return db.Jalab{}, fmt.Errorf("executing query: %w", errCreate)
	}

	return created, nil
}

func (r Repo) Get(ctx context.Context, j db.Jalab) (gotten db.Jalab, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("getting jalab %+v: %w", j, err)
	}()

	errGet := r.stmtGet.GetContext(ctx, &gotten, j)
	if errors.Is(errGet, sql.ErrNoRows) {
		return db.Jalab{}, ErrNotFound
	}
	if errGet != nil {
		return db.Jalab{}, fmt.Errorf("executing query: %w", errGet)
	}

	return gotten, nil
}

func (r Repo) GetAll(ctx context.Context, j db.Jalab) (gotten []db.Jalab, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("getting all jalabs of group with %d: %w", j.GroupChatID, err)
	}()

	errGetAll := r.stmtGetAll.SelectContext(ctx, &gotten, j)
	if errors.Is(errGetAll, sql.ErrNoRows) {
		return []db.Jalab{}, nil
	}
	if errGetAll != nil {
		return nil, fmt.Errorf("executing query: %w", errGetAll)
	}

	return gotten, nil
}
