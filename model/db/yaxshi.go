package db

import (
	"time"
)

type Yaxshi struct {
	ID          int64     `db:"id"`
	Count       int64     `db:"count"`
	GroupChatID int64     `db:"group_chat_id"`
	UserID      int64     `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
