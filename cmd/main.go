package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"jalabs.kz/bot/model"
	"jalabs.kz/bot/service"
	"jalabs.kz/bot/storage"
	"jalabs.kz/bot/tgapi"
	"jalabs.kz/bot/tgapi/command_scope"
)

const (
	becomeJalab    = "becomejalab"
	todaysJalab    = "todaysjalab"
	jalabs         = "jalabs"
	yaxshi         = "yaxshi"
	toggleMentions = "togglementions"

	jalabsOfTheMonth = "jalabsofthemonth"
)

func main() {
	tokenPtr := flag.String("bot_token", "", "Telegram API access token for the bot")
	flag.Parse()
	if tokenPtr == nil {
		log.Fatalln("failed to parse bot_token flag")
	}

	botToken := *tokenPtr

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stg, errNewStorage := storage.New(ctx)
	if errNewStorage != nil {
		log.Fatalln(errNewStorage)
	}

	svc := service.New(stg)

	go func() {
		sig := <-sigChan
		log.Printf("received signal %v\n", sig)
		cancel()
	}()

	tg := tgapi.New(botToken)

	tg.RegisterCommand(tgapi.CommandConfig{
		Command:   "help",
		Scopes:    []string{command_scope.Private, command_scope.Groups},
		ShortDesc: "Подробное описание команд",
		Handler:   Help,
	})
	tg.RegisterCommand(tgapi.CommandConfig{
		Command:   becomeJalab,
		Scopes:    []string{command_scope.Groups},
		ShortDesc: "Стать джаляпом",
		Handler:   svc.BecomeJalab,
	})
	tg.RegisterCommand(tgapi.CommandConfig{
		Command:   todaysJalab,
		Scopes:    []string{command_scope.Groups},
		ShortDesc: "Главный джаляп дня",
		Handler:   svc.TodaysJalab,
	})
	tg.RegisterCommand(tgapi.CommandConfig{
		Command:   jalabs,
		Scopes:    []string{command_scope.Groups},
		ShortDesc: "Список джаляпов",
		Handler:   svc.Jalabs,
	})
	tg.RegisterCommand(tgapi.CommandConfig{
		Command:   yaxshi,
		Scopes:    []string{command_scope.Groups},
		ShortDesc: "Статистика яхши / Поблагодарить",
		Handler:   svc.Yaxshi,
	})
	tg.RegisterCommand(tgapi.CommandConfig{
		Command:   toggleMentions,
		Scopes:    []string{command_scope.Private, command_scope.Groups},
		ShortDesc: "Вкл/выкл свои упоминания",
		Handler:   svc.ToggleMentions,
	})
	tg.RegisterCommand(tgapi.CommandConfig{
		Command:   jalabsOfTheMonth,
		Scopes:    []string{command_scope.Groups},
		ShortDesc: "Джаляпы месяца",
		Handler:   svc.JalabsOfTheMonth,
	})
	errStart := tg.Start(ctx)
	if errStart != nil {
		log.Println(errStart)
	}
}

func Help(c tgapi.HandlerContext, u model.Update, _ ...string) error {
	log.Printf("received /help command from %q\n", u.Message.From.Username)

	return c.SendMessage(model.SendMessageRq{
		ChatID: u.Message.Chat.Id,
		Text:   "НЕГР",
	})
}
