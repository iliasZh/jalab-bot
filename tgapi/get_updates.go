package tgapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"path"

	"golang.org/x/exp/slices"
	"jalabs.kz/bot/model"
)

func (t *TgAPI) getUpdatesURL() *url.URL {
	botToken := fmt.Sprintf("bot%s", t.botToken)
	getUpdatesURL := &url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   path.Join(botToken, "getUpdates"),
	}

	queryParams := url.Values{}
	queryParams.Set("offset", fmt.Sprintf("%v", t.maxUpdateID+1))
	queryParams.Set("timeout", fmt.Sprintf("%d", pollingPeriodSeconds))

	getUpdatesURL.RawQuery = queryParams.Encode()
	return getUpdatesURL
}

func (t *TgAPI) getUpdates(ctx context.Context, updatesChan chan<- model.Update) {
	for {
		select {
		case <-ctx.Done():
			close(updatesChan)
			return
		default:
			updates, errDo := doRequest[model.GetUpdatesRq, model.GetUpdatesRs](
				ctx, t, "getUpdates",
				model.GetUpdatesRq{
					Offset:  t.maxUpdateID + 1,
					Timeout: pollingPeriodSeconds,
				},
			)
			if errors.Is(errDo, context.Canceled) {
				continue
			}
			if errDo != nil {
				log.Println(errDo)
				continue
			}

			slices.SortFunc(updates, func(u1, u2 model.Update) bool {
				return u1.Message.Date < u2.Message.Date
			})

			if updatesCount := len(updates); updatesCount != 0 {
				t.maxUpdateID = updates[updatesCount-1].Id
			}

			for _, u := range updates {
				updatesChan <- u
			}
		}
	}
}
