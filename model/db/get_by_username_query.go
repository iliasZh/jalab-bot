package db

type GetByUsernameQuery struct {
	Username    string `db:"username"`
	GroupChatID int64  `db:"group_chat_id"`
}
