package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"jalabs.kz/bot/model"
	"jalabs.kz/bot/model/db"
	"jalabs.kz/bot/storage/postgres"
	"jalabs.kz/bot/tgapi"
)

func (s Service) Yaxshi(c tgapi.HandlerContext, u model.Update, args ...string) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	c.SetContext(ctx)

	if len(args) != 0 {
		return s.giveYaxshi(c, u, args[0])
	}

	return s.getYaxshiStats(c, u)
}

func (s Service) giveYaxshi(c tgapi.HandlerContext, u model.Update, username string) error {
	username, ok := strings.CutPrefix(username, "@")
	if !ok {
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     fmt.Sprintf("%q не юзернейм!", username),
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	if username == c.Bot().Username {
		text := "Рахмет, брат"
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     text,
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	if u.Message.From.Username == username {
		text := fmt.Sprintf("%s, меня не обмануть!", u.Message.From.FirstName)
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     text,
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	jalab, errGetReceiver := s.stg.Repo.GetJalabByUsername(c.Context(), db.GetByUsernameQuery{
		Username:    username,
		GroupChatID: u.Message.Chat.Id,
	})
	if errors.Is(errGetReceiver, postgres.ErrJalabNotFound) {
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     fmt.Sprintf("Джаляп @%s не найден!", username),
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}
	if errGetReceiver != nil {
		return errGetReceiver
	}

	yaxshi, errCreateYaxshi := s.stg.Yaxshis.Create(c.Context(), db.Yaxshi{
		GroupChatID: u.Message.Chat.Id,
		UserID:      jalab.ID,
	})
	if errCreateYaxshi != nil {
		return errCreateYaxshi
	}

	text := fmt.Sprintf(
		"Яхши %s %s! Всего %d яхши", jalab.FirstName, jalab.UsernameString(), yaxshi.Count,
	)
	return c.SendMessage(model.SendMessageRq{
		ChatID:                   u.Message.Chat.Id,
		Text:                     text,
		ReplyToMessageID:         u.Message.Id,
		AllowSendingWithoutReply: true,
	})
}

func (s Service) getYaxshiStats(c tgapi.HandlerContext, u model.Update) error {
	jalabs, stats, errGetStats := s.stg.Repo.GetYaxshiStats(c.Context(), db.Yaxshi{
		GroupChatID: u.Message.Chat.Id,
	})
	if errGetStats != nil {
		return errGetStats
	}

	numEntries := len(jalabs)
	if numEntries == 0 {
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     "Нету яхши!",
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	entries := make([]string, numEntries)
	for i := range entries {
		j := jalabs[i]
		entries[i] = fmt.Sprintf(
			"%d. %s %s - %d яхши", i+1, j.FirstName, j.UsernameString(), stats[i],
		)
	}

	text := fmt.Sprintf("Статистика яхши:\n") + strings.Join(entries, "\n")
	return c.SendMessage(model.SendMessageRq{
		ChatID:                   u.Message.Chat.Id,
		Text:                     text,
		ReplyToMessageID:         u.Message.Id,
		AllowSendingWithoutReply: true,
	})
}
