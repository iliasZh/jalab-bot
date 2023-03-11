package db

import (
	"time"
)

type Jalab struct {
	ID          int64     `db:"id"`
	GroupChatID int64     `db:"group_chat_id"`
	UserID      int64     `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
