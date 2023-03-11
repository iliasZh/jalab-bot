package model

type Update struct {
	Id      int64   `json:"update_id"`
	Message Message `json:"message"`
}
