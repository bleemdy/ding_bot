package ding_bot

import (
	"github.com/bleemdy/ding_bot/context"
	"github.com/bleemdy/ding_bot/message"
	"github.com/bleemdy/ding_bot/schedule"
)

func New() *context.Bot {
	botManage := &context.Bot{
		Messages:       make(chan *context.Context, 10),
		Commands:       map[string]func(ctx *context.Context){},
		MessageManager: message.New(),
		Schedule:       schedule.New(),
	}
	return botManage
}
