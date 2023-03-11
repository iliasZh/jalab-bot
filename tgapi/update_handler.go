package tgapi

import (
	"fmt"
	"strings"

	"jalabs.kz/bot/model"
	"jalabs.kz/bot/tgapi/chat_type"
)

type Handler func(c HandlerContext, u model.Update, optionalArgs ...string) error

type CommandConfig struct {
	Command   string
	Scopes    []string
	ShortDesc string
	Handler   Handler
}

func (t *TgAPI) RegisterCommand(cfg CommandConfig) {
	if len(cfg.Scopes) == 0 {
		panic("no scopes provided")
	}

	if cfg.Command == "" {
		panic("cannot register empty command")
	}

	if cfg.Handler == nil {
		panic("handler is nil")
	}

	for _, scope := range cfg.Scopes {
		key := handlerKey(scope, cfg.Command)

		_, ok := t.handlers[key]
		if ok {
			panic(fmt.Sprintf("handler for key %q is already registered", cfg.Command))
		}

		t.handlers[key] = cfg.Handler

		t.commands[scope] = append(t.commands[scope], model.BotCommand{
			Command:     cfg.Command,
			Description: cfg.ShortDesc,
		})
	}
}

func (t *TgAPI) handleUpdate(u model.Update) error {
	c := t.newContext()

	text := u.Message.Text
	if !strings.HasPrefix(text, "/") {
		return nil
	}

	parts := strings.Split(text, " ")
	noPrefix := strings.TrimPrefix(parts[0], "/")

	cmd := strings.TrimSuffix(noPrefix, fmt.Sprintf("@%s", t.self.Username))
	scope := chat_type.ToScope(u.Message.Chat.Type)

	key := handlerKey(scope, cmd)

	hnd, ok := t.handlers[key]
	if !ok {
		return nil
	}

	args := parts[1:]
	return hnd(c, u, args...)
}

func handlerKey(scope, cmd string) string {
	return fmt.Sprintf("%s/%s", scope, cmd)
}
