package context

import (
	"github.com/bleemdy/ding_bot/context/internal"
	"github.com/bleemdy/ding_bot/message"
	"strings"
)

type Context struct {
	bot     *Bot
	Message *message.Ding
	Content string
	Webhook string
	Args    []string
}

func (c Context) SendText(text string) {
	content := message.Text{
		Basic: message.Basic{
			Text:      text,
			Webhook:   c.Webhook,
			AtUserIds: c.Message.SenderStaffId,
		},
	}
	c.bot.messageManager.Send(content)
}

func (c Context) SendMarkDown(title, text string) {
	content := message.Markdown{
		Title: title,
		Basic: message.Basic{
			Text:      text,
			Webhook:   c.Webhook,
			AtUserIds: c.Message.SenderStaffId,
		},
	}
	c.bot.messageManager.Send(content)
}

func (c Context) SendActionCard(title, text, singleTitle, singleURL string) {
	content := message.ActionCard{
		Title:       title,
		SingleURL:   singleURL,
		SingleTitle: singleTitle,
		Basic: message.Basic{
			Text:      text,
			Webhook:   c.Webhook,
			AtUserIds: c.Message.SenderStaffId,
		},
	}
	c.bot.messageManager.Send(content)
}

func (c Context) Send(msg message.Common) {
	c.bot.messageManager.Send(msg)
}

func newContext(b *Bot, d *message.Ding) *Context {
	args := strings.Split(strings.TrimSpace(d.Text.Content), " ")
	content := internal.CompressStr(d.Text.Content, "")
	return &Context{
		b,
		d,
		content,
		d.SessionWebhook,
		args,
	}
}
