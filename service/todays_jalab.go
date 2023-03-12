package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"jalabs.kz/bot/model"
	"jalabs.kz/bot/model/db"
	"jalabs.kz/bot/storage/postgres"
	"jalabs.kz/bot/tgapi"
)

func (s Service) TodaysJalab(c tgapi.HandlerContext, u model.Update, _ ...string) error {
	ctx, cancel := context.WithTimeout(c.Ctx(), 5*time.Second)
	defer cancel()

	c.SetCtx(ctx)

	jalabUser, errGet := s.stg.Repo.GetTodaysJalab(c.Ctx(), db.TodaysJalab{
		GroupChatID: u.Message.Chat.Id,
	})
	if errGet == nil {
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     todaysJalabText(jalabUser),
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}
	if errGet != nil && !errors.Is(errGet, postgres.ErrTodaysJalabNotFound) {
		return errGet
	}

	groupJalabs, errGetGroup := s.stg.Jalabs.GetAll(c.Ctx(), db.Jalab{GroupChatID: u.Message.Chat.Id})
	if errGetGroup != nil {
		return errGetGroup
	}

	jalabsCount := len(groupJalabs)
	if jalabsCount == 0 {
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     "Не добавлено ни одного джаляпа!",
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	idx := rand.Intn(jalabsCount)
	jalab := groupJalabs[idx]
	todaysJalab, errCreateOrGet := s.stg.TodaysJalabs.CreateOrGet(c.Ctx(), db.TodaysJalab{
		UserID:      jalab.UserID,
		GroupChatID: jalab.GroupChatID,
	})
	if errCreateOrGet != nil {
		return errCreateOrGet
	}

	jalabUser, errGetUser := s.stg.Users.Get(c.Ctx(), db.User{ID: todaysJalab.UserID})
	if errGetUser != nil {
		return errGetUser
	}
	return c.SendMessage(model.SendMessageRq{
		ChatID:                   u.Message.Chat.Id,
		Text:                     todaysJalabText(jalabUser),
		ReplyToMessageID:         u.Message.Id,
		AllowSendingWithoutReply: true,
	})
}

func todaysJalabText(u db.User) string {
	return fmt.Sprintf("Главный джаляп дня - %s (%s)", u.FirstName, u.UsernameString())
}
