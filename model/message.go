package model

type Message struct {
	Id       int64    `json:"message_id"`
	From     User     `json:"from"`
	Chat     Chat     `json:"chat"`
	Date     int      `json:"date"`
	Text     string   `json:"text"`
	Entities []Entity `json:"entities"`
}
