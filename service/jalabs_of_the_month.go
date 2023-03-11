package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"jalabs.kz/bot/model"
	"jalabs.kz/bot/model/db"
	"jalabs.kz/bot/tgapi"
)

func (s Service) JalabsOfTheMonth(c tgapi.HandlerContext, u model.Update, _ ...string) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	c.SetContext(ctx)

	q := db.NewGetJalabOfTheMonthQuery(u.Message.Chat.Id)
	jalabs, stats, errGet := s.stg.Repo.GetJalabStats(c.Context(), q)
	if errGet != nil {
		return errGet
	}

	numEntries := len(jalabs)
	if numEntries == 0 {
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     "Джаляпов дня ещё не было!",
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	entries := make([]string, numEntries)
	for i := range entries {
		j := jalabs[i]
		entries[i] = fmt.Sprintf(
			"%d. %s %s - %dx джаляп", i+1, j.FirstName, j.UsernameString(), stats[i],
		)
	}

	text := fmt.Sprintf("Джаляпы месяца:\n") + strings.Join(entries, "\n")
	return c.SendMessage(model.SendMessageRq{
		ChatID:                   u.Message.Chat.Id,
		Text:                     text,
		ReplyToMessageID:         u.Message.Id,
		AllowSendingWithoutReply: true,
	})
}
