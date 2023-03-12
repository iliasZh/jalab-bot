package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"jalabs.kz/bot/model/db"
)

var (
	ErrJalabNotFound       = fmt.Errorf("jalab not found")
	ErrTodaysJalabNotFound = fmt.Errorf("today's jalab not found")
)

const (
	queryGetTodaysJalab = `
	SELECT u.id, u.username, u.first_name, u.disable_mentions, u.created_at, u.updated_at
	FROM todays_jalabs tj
	    LEFT JOIN jalabs j ON tj.group_chat_id = j.group_chat_id AND tj.user_id = j.user_id
		LEFT JOIN users u ON j.user_id = u.id
	WHERE tj.group_chat_id = :group_chat_id AND tj.created_at = CURRENT_DATE;`

	queryGetGroupJalabs = `
	SELECT u.id, u.username, u.first_name, u.disable_mentions, u.created_at, u.updated_at
	FROM jalabs j
	    LEFT JOIN users u ON j.user_id = u.id
	WHERE group_chat_id = :group_chat_id;`

	queryGetJalabByUsername = `
	SELECT u.id, u.username, u.first_name, u.disable_mentions, u.created_at, u.updated_at
	FROM jalabs j
	    LEFT JOIN users u ON j.user_id = u.id
	WHERE j.group_chat_id = :group_chat_id AND u.username = :username;`

	queryGetYaxshiStats = `
	SELECT u.id, u.username, u.first_name, u.disable_mentions, u.created_at, u.updated_at, y.count
	FROM yaxshis y
	    LEFT JOIN jalabs j ON j.group_chat_id = y.group_chat_id AND j.user_id = y.user_id
	    LEFT JOIN users u ON u.id = j.user_id
	WHERE y.group_chat_id = :group_chat_id
	ORDER BY count DESC, y.updated_at DESC;`

	queryGetJalabStats = `
	SELECT u.id, u.username, u.first_name, u.disable_mentions,
	       COUNT(tj.id) AS jalab_count, u.created_at, u.updated_at
	FROM todays_jalabs tj
    	LEFT JOIN jalabs j ON tj.group_chat_id = j.group_chat_id AND tj.user_id = j.user_id
    	LEFT JOIN users u ON j.user_id = u.id
	WHERE tj.group_chat_id = :group_chat_id
		AND tj.created_at BETWEEN :from_date AND :to_date
	GROUP BY u.id
	ORDER BY jalab_count DESC;`
)

type Repo struct {
	db *sqlx.DB

	stmtGetTodaysJalab     *sqlx.NamedStmt
	stmtGetGroupJalabs     *sqlx.NamedStmt
	stmtGetJalabByUsername *sqlx.NamedStmt
	stmtGetYaxshiStats     *sqlx.NamedStmt
	stmtGetJalabStats      *sqlx.NamedStmt
}

func NewRepo(ctx context.Context, db *sqlx.DB) (r Repo, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("new repo: %w", err)
	}()

	stmtGetTodaysJalab, errPrepareGet := db.PrepareNamedContext(ctx, queryGetTodaysJalab)
	if errPrepareGet != nil {
		return Repo{}, fmt.Errorf("preparing get todays jalab stmt: %w", errPrepareGet)
	}

	stmtGetGroupJalabs, errPrepareGetGroupJalabs := db.PrepareNamedContext(ctx, queryGetGroupJalabs)
	if errPrepareGetGroupJalabs != nil {
		return Repo{}, fmt.Errorf(
			"preparing get group jalabs stmt: %w", errPrepareGetGroupJalabs,
		)
	}

	stmtGetJalabByUsername, errPrepareGetJalabByUsername := db.PrepareNamedContext(
		ctx, queryGetJalabByUsername,
	)
	if errPrepareGetJalabByUsername != nil {
		return Repo{}, fmt.Errorf(
			"preparing get jalab by username stmt: %w", errPrepareGetJalabByUsername,
		)
	}

	stmtGetYaxshiStats, errPrepareGetYaxshiStats := db.PrepareNamedContext(ctx, queryGetYaxshiStats)
	if errPrepareGetYaxshiStats != nil {
		return Repo{}, fmt.Errorf(
			"preparing get yaxshi stats stmt: %w", errPrepareGetYaxshiStats,
		)
	}

	stmtGetJalabStats, errPrepareGetJalabStats := db.PrepareNamedContext(ctx, queryGetJalabStats)
	if errPrepareGetJalabStats != nil {
		return Repo{}, fmt.Errorf(
			"preparing get jalab stats stmt: %w", errPrepareGetJalabStats,
		)
	}

	return Repo{
		db:                     db,
		stmtGetTodaysJalab:     stmtGetTodaysJalab,
		stmtGetGroupJalabs:     stmtGetGroupJalabs,
		stmtGetJalabByUsername: stmtGetJalabByUsername,
		stmtGetYaxshiStats:     stmtGetYaxshiStats,
		stmtGetJalabStats:      stmtGetJalabStats,
	}, nil
}

func (r Repo) GetTodaysJalab(
	ctx context.Context,
	tj db.TodaysJalab,
) (u db.User, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("getting todays jalab: %w", err)
	}()

	errGet := r.stmtGetTodaysJalab.GetContext(ctx, &u, tj)
	if errors.Is(errGet, sql.ErrNoRows) {
		return db.User{}, ErrTodaysJalabNotFound
	}
	if errGet != nil {
		return db.User{}, fmt.Errorf("executing query: %w", errGet)
	}

	return u, nil
}

func (r Repo) GetGroupJalabs(ctx context.Context, j db.Jalab) (users []db.User, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("getting group jalabs: %w", err)
	}()

	errGet := r.stmtGetGroupJalabs.SelectContext(ctx, &users, j)
	if errGet != nil {
		return nil, fmt.Errorf("executing query: %w", errGet)
	}

	return users, nil
}

func (r Repo) GetJalabByUsername(
	ctx context.Context,
	q db.GetByUsernameQuery,
) (u db.User, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("getting jalab by username %+v: %w", q, err)
	}()

	errGet := r.stmtGetJalabByUsername.GetContext(ctx, &u, q)
	if errors.Is(errGet, sql.ErrNoRows) {
		return db.User{}, ErrJalabNotFound
	}
	if errGet != nil {
		return db.User{}, fmt.Errorf("executing query: %w", errGet)
	}

	return u, nil
}

func (r Repo) GetYaxshiStats(
	ctx context.Context,
	y db.Yaxshi,
) (jalabs []db.User, stats []int64, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("getting yaxshi stats for group %d: %w", y.GroupChatID, err)
	}()

	result := make([]struct {
		db.User
		db.Yaxshi
	}, 0)

	errGet := r.stmtGetYaxshiStats.SelectContext(ctx, &result, y)
	if errors.Is(errGet, sql.ErrNoRows) {
		return nil, nil, nil
	}
	if errGet != nil {
		return nil, nil, errGet
	}

	jalabs = make([]db.User, len(result))
	stats = make([]int64, len(result))
	for i, r := range result {
		jalabs[i] = r.User
		stats[i] = r.Yaxshi.Count
	}

	return jalabs, stats, nil
}

func (r Repo) GetJalabStats(
	ctx context.Context,
	q db.GetJalabStatsQuery,
) (jalabs []db.User, stats []int64, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("getting jalab stats for query %+v: %w", q, err)
	}()

	results := make([]struct {
		db.User
		JalabCount int64 `db:"jalab_count"`
	}, 0)

	errGet := r.stmtGetJalabStats.SelectContext(ctx, &results, q)
	if errors.Is(errGet, sql.ErrNoRows) {
		return nil, nil, nil
	}
	if errGet != nil {
		return nil, nil, fmt.Errorf("executing query: %w", errGet)
	}

	jalabs = make([]db.User, len(results))
	stats = make([]int64, len(results))
	for i, result := range results {
		jalabs[i] = result.User
		stats[i] = result.JalabCount
	}

	return jalabs, stats, nil
}
