package db

import (
	"time"
)

type TodaysJalab struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	GroupChatID int64     `db:"group_chat_id"`
	CreatedAt   time.Time `db:"created_at"`
}
