package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	"jalabs.kz/bot/model"
	"jalabs.kz/bot/model/db"
	"jalabs.kz/bot/tgapi"
)

func (s Service) Jalabs(c tgapi.HandlerContext, u model.Update, _ ...string) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	c.SetContext(ctx)

	users, errGet := s.stg.Repo.GetGroupJalabs(
		c.Context(), db.Jalab{GroupChatID: u.Message.Chat.Id},
	)
	if errGet != nil {
		return errGet
	}

	usersCount := len(users)
	if usersCount == 0 {
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     "Нет джаляпов!",
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	slices.SortFunc(users, func(u1, u2 db.User) bool {
		return u1.CreatedAt.Before(u2.CreatedAt)
	})

	jalabList := make([]string, len(users))
	for i, jalab := range users {
		jalabList[i] = fmt.Sprintf("%d. %s (%s)", i+1, jalab.FirstName, jalab.UsernameString())
	}

	return c.SendMessage(model.SendMessageRq{
		ChatID:                   u.Message.Chat.Id,
		Text:                     "Список джаляпов\n" + strings.Join(jalabList, "\n"),
		ReplyToMessageID:         u.Message.Id,
		AllowSendingWithoutReply: true,
	})
}
