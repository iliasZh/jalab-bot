package tgapi

import (
	"context"
	"log"
	"net/http"
	"time"

	"jalabs.kz/bot/model"
)

const (
	pollingPeriodSeconds   = 50
	requestOverheadSeconds = 2
)

type TgAPI struct {
	cli         *http.Client
	commands    map[string][]model.BotCommand
	handlers    map[string]Handler
	botToken    string
	maxUpdateID int64
	self        model.User
}

func New(botToken string) TgAPI {
	return TgAPI{
		cli: &http.Client{
			Timeout: time.Duration(pollingPeriodSeconds+requestOverheadSeconds) * time.Second,
		},
		commands:    make(map[string][]model.BotCommand),
		handlers:    make(map[string]Handler),
		botToken:    botToken,
		maxUpdateID: 0,
	}
}

func (t *TgAPI) Start(ctx context.Context) error {
	self, errGetSelf := t.Self(ctx)
	if errGetSelf != nil {
		return errGetSelf
	}

	t.self = self

	for scope, cmd := range t.commands {
		_, errDo := doRequest[model.SetMyCommandsRq, bool](
			ctx, t, "setMyCommands",
			model.SetMyCommandsRq{
				Commands: cmd,
				Scope: model.BotCommandScope{
					Type: scope,
				},
			},
		)
		if errDo != nil {
			return errDo
		}
	}

	updatesChan := make(chan model.Update, 10)

	go t.getUpdates(ctx, updatesChan)

	for u := range updatesChan {
		log.Printf("%+v\n", u)

		go func(u model.Update) {
			errHandle := t.handleUpdate(u)
			if errHandle != nil {
				log.Println(errHandle)
			}
		}(u)
	}

	return nil
}

func (t *TgAPI) Self(ctx context.Context) (model.User, error) {
	self, errDo := doRequest[any, model.User](ctx, t, "getMe", nil)
	if errDo != nil {
		return model.User{}, errDo
	}

	return self, nil
}

func (t *TgAPI) newContext() HandlerContext {
	return HandlerContext{
		ctx: context.Background(),
		t:   t,
	}
}
