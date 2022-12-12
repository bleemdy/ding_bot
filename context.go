package ding_bot

import (
	"github.com/bleemdy/ding_bot/push"
	"math"
	"net/http"
	"strings"
)

const abortIndex int8 = math.MaxInt8 >> 1

type Context struct {
	index   int8
	Bot     *Bot
	Message *push.Ding
	// 消息内容
	Content string
	Webhook string
	Args    []string
	Command string
	// 1：单聊
	// 2：群聊
	ConversationType string
	IsAdmin          bool
	// 发送者ID
	SenderStaffId  string
	handlers       []func(*Context)
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Run() {
	if len(c.handlers) == 0 {
		return
	}
	c.handlers[0](c)
}

func (c *Context) SendText(text string) {
	content := push.Text{
		Text:      text,
		Webhook:   c.Webhook,
		AtUserIds: c.Message.SenderStaffId,
	}
	c.Bot.Send(content)
}

func (c *Context) SendMarkDown(title, text string) {
	content := push.Markdown{
		Title:     title,
		Text:      text,
		Webhook:   c.Webhook,
		AtUserIds: c.Message.SenderStaffId,
	}
	c.Bot.Send(content)
}

func (c *Context) SendActionCard(title, text, singleTitle, singleURL string) {
	content := push.ActionCard{
		Title:       title,
		SingleURL:   singleURL,
		SingleTitle: singleTitle,
		Text:        text,
		Webhook:     c.Webhook,
		AtUserIds:   c.Message.SenderStaffId,
	}
	c.Bot.Send(content)
}

func (c *Context) Send(msg push.Common) {
	c.Bot.Send(msg)
}

func newContext(ding *push.Ding) *Context {
	isCommand := strings.Contains(ding.Text.Content, "/")
	var args []string
	var command string
	content := strings.TrimSpace(ding.Text.Content)
	if isCommand {
		args = strings.Split(content, " ")
		if len(args) >= 1 {
			command = args[0]
			args = args[1:]
			content = strings.TrimSpace(strings.TrimLeft(content, command))
		} else {
			return nil
		}
	}
	return &Context{
		Message:          ding,
		Content:          content,
		Args:             args,
		Command:          command,
		Webhook:          ding.SessionWebhook,
		ConversationType: ding.ConversationType,
		IsAdmin:          ding.IsAdmin,
	}
}
