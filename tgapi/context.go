package tgapi

import (
	"context"

	"jalabs.kz/bot/model"
)

type HandlerContext struct {
	ctx context.Context
	t   *TgAPI
}

func (c *HandlerContext) SetCtx(ctx context.Context) {
	c.ctx = ctx
}

func (c *HandlerContext) Ctx() context.Context {
	return c.ctx
}

func (c *HandlerContext) Bot() model.User {
	return c.t.self
}

func (c *HandlerContext) SendMessage(rq model.SendMessageRq) (err error) {
	_, errDo := doRequest[model.SendMessageRq, model.SendMessageRs](
		c.ctx, c.t, "sendMessage", rq,
	)
	return errDo
}
