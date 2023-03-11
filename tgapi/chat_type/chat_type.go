package chat_type

import (
	"jalabs.kz/bot/tgapi/command_scope"
)

const (
	Private    = "private"
	Group      = "group"
	Supergroup = "supergroup"
)

func ToScope(chatType string) string {
	switch chatType {
	case Private:
		return command_scope.Private
	case Group, Supergroup:
		return command_scope.Groups
	default:
		return ""
	}
}
