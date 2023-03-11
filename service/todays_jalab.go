package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"jalabs.kz/bot/model"
	"jalabs.kz/bot/model/db"
	"jalabs.kz/bot/tgapi"
)

func (s Service) TodaysJalab(c tgapi.HandlerContext, u model.Update, _ ...string) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	c.SetContext(ctx)

	jalabUser, errGet := s.stg.Repo.GetTodaysJalab(ctx, db.TodaysJalab{
		GroupChatID: u.Message.Chat.Id,
	})
	if errGet != nil && !errors.Is(errGet, sql.ErrNoRows) {
		return errGet
	}
	if errGet == nil {
		text := fmt.Sprintf(
			"Главный джаляп дня - %s (%s)", jalabUser.FirstName, jalabUser.UsernameString(),
		)
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     text,
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	members, errGetGroup := s.stg.Repo.GetGroupJalabs(ctx, db.Jalab{GroupChatID: u.Message.Chat.Id})
	if errGetGroup != nil {
		return errGetGroup
	}

	jalabsCount := len(members)
	if jalabsCount == 0 {
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     "Не добавлено ни одного джаляпа!",
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	idx := rand.Intn(jalabsCount)
	jalab := members[idx]

	_, errCreate := s.stg.TodaysJalabs.Create(ctx, db.TodaysJalab{
		UserID:      jalab.ID,
		GroupChatID: u.Message.Chat.Id,
	})
	if errCreate != nil {
		return errCreate
	}

	text := fmt.Sprintf(
		"Главный джаляп дня - %s (%s)", jalab.FirstName, jalab.UsernameString(),
	)
	return c.SendMessage(model.SendMessageRq{
		ChatID:                   u.Message.Chat.Id,
		Text:                     text,
		ReplyToMessageID:         u.Message.Id,
		AllowSendingWithoutReply: true,
	})
}
