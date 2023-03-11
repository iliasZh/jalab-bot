package model

type SetMyCommandsRq struct {
	Commands []BotCommand    `json:"commands"`
	Scope    BotCommandScope `json:"scope"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type BotCommandScope struct {
	Type string `json:"type"`
}
