package ding_bot

import (
	"github.com/bleemdy/ding_bot/message"
	"github.com/bleemdy/ding_bot/utils"
	"strings"
)

type Context struct {
	Bot     *Bot
	Message *message.Ding
	Content string
	Webhook string
	Args    []string
	Command string
}

func (c Context) SendText(text string) {
	content := message.Text{
		Text:      text,
		Webhook:   c.Webhook,
		AtUserIds: c.Message.SenderStaffId,
	}
	c.Bot.MessageManager.Send(content)
}

func (c Context) SendMarkDown(title, text string) {
	content := message.Markdown{
		Title:     title,
		Text:      text,
		Webhook:   c.Webhook,
		AtUserIds: c.Message.SenderStaffId,
	}
	c.Bot.MessageManager.Send(content)
}

func (c Context) SendActionCard(title, text, singleTitle, singleURL string) {
	content := message.ActionCard{
		Title:       title,
		SingleURL:   singleURL,
		SingleTitle: singleTitle,
		Text:        text,
		Webhook:     c.Webhook,
		AtUserIds:   c.Message.SenderStaffId,
	}
	c.Bot.MessageManager.Send(content)
}

func (c Context) Send(msg message.Common) {
	c.Bot.MessageManager.Send(msg)
}

func newContext(b *Bot, d *message.Ding) *Context {
	args := strings.Split(strings.TrimSpace(d.Text.Content), " ")
	command := ""
	if len(args) > 0 {
		command = args[0]
	}
	content := utils.CompressStr(d.Text.Content, "")
	return &Context{
		Bot:     b,
		Message: d,
		Content: content,
		Args:    args[1:],
		Command: command,
		Webhook: d.SessionWebhook,
	}
}
