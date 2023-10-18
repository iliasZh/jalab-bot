package tgapi

import (
	"cmp"
	"context"
	"errors"
	"log"
	"slices"

	"jalabs.kz/bot/model"
)

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

			slices.SortFunc(updates, func(u1, u2 model.Update) int {
				return cmp.Compare(u1.Message.Date, u2.Message.Date)
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
