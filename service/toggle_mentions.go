package service

import (
	"context"
	"errors"
	"time"

	"jalabs.kz/bot/model"
	"jalabs.kz/bot/model/db"
	"jalabs.kz/bot/storage/postgres/users"
	"jalabs.kz/bot/tgapi"
)

func (s Service) ToggleMentions(c tgapi.HandlerContext, u model.Update, _ ...string) error {
	ctx, cancel := context.WithTimeout(c.Ctx(), 5*time.Second)
	defer cancel()

	c.SetCtx(ctx)

	user, errToggle := s.stg.Users.ToggleMentions(c.Ctx(), db.User{
		ID: u.Message.From.Id,
	})
	if errors.Is(errToggle, users.ErrNotFound) {
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     "Тебя нету в базе! Я пока не могу тебя упоминать",
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	if errToggle != nil {
		return errToggle
	}

	text := "Упоминания включены!"
	if user.DisableMentions {
		text = "Упоминания отключены!"
	}

	return c.SendMessage(model.SendMessageRq{
		ChatID:                   u.Message.Chat.Id,
		Text:                     text,
		ReplyToMessageID:         u.Message.Id,
		AllowSendingWithoutReply: true,
	})
}
