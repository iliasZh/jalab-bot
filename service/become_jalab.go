package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"jalabs.kz/bot/model"
	"jalabs.kz/bot/model/db"
	"jalabs.kz/bot/storage/postgres/jalabs"
	"jalabs.kz/bot/tgapi"
)

func (s Service) BecomeJalab(c tgapi.HandlerContext, u model.Update, _ ...string) error {
	ctx, cancel := context.WithTimeout(c.Ctx(), 5*time.Second)
	defer cancel()
	c.SetCtx(ctx)

	storedUser, errGetUser := s.stg.Users.Get(c.Ctx(), db.User{ID: u.Message.From.Id})
	if errGetUser != nil && !errors.Is(errGetUser, sql.ErrNoRows) {
		return errGetUser
	}

	currentUser := db.User{
		ID: u.Message.From.Id,
		Username: sql.NullString{
			String: u.Message.From.Username,
			Valid:  u.Message.From.Username != "",
		},
		FirstName: u.Message.From.FirstName,
	}

	currentUser, errCreateUser := s.stg.Users.Create(c.Ctx(), currentUser)
	if errCreateUser != nil {
		return errCreateUser
	}

	if errGetUser == nil {
		updateText := make([]string, 0, 1)
		if storedUser.Username != currentUser.Username {
			updateText = append(updateText, fmt.Sprintf(
				"Юзернейм: %s -> %s", storedUser.UsernameString(), currentUser.UsernameString(),
			))
		}
		if storedUser.FirstName != currentUser.FirstName {
			updateText = append(updateText, fmt.Sprintf(
				"Имя: %s -> %s", storedUser.FirstName, currentUser.FirstName,
			))
		}
		if len(updateText) != 0 {
			text := "Джаляп обновлён!\n" + strings.Join(updateText, "\n")
			return c.SendMessage(model.SendMessageRq{
				ChatID:                   u.Message.Chat.Id,
				Text:                     text,
				ReplyToMessageID:         u.Message.Id,
				AllowSendingWithoutReply: true,
			})
		}
	}

	jalab := db.Jalab{
		GroupChatID: u.Message.Chat.Id,
		UserID:      currentUser.ID,
	}
	_, errCreate := s.stg.Jalabs.Create(c.Ctx(), jalab)
	if errCreate != nil && !errors.Is(errCreate, jalabs.ErrAlreadyExists) {
		return errCreate
	}
	if errCreate == nil {
		text := fmt.Sprintf("Джаляп %s (%s) добавлен!", currentUser.FirstName, currentUser.UsernameString())
		return c.SendMessage(model.SendMessageRq{
			ChatID:                   u.Message.Chat.Id,
			Text:                     text,
			ReplyToMessageID:         u.Message.Id,
			AllowSendingWithoutReply: true,
		})
	}

	text := fmt.Sprintf("Джаляп %s (%s) уже был добавлен", currentUser.FirstName, currentUser.UsernameString())
	return c.SendMessage(model.SendMessageRq{
		ChatID:                   u.Message.Chat.Id,
		Text:                     text,
		ReplyToMessageID:         u.Message.Id,
		AllowSendingWithoutReply: true,
	})
}
