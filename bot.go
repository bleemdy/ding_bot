package ding_bot

import (
	"encoding/json"
	"fmt"
	"github.com/bleemdy/ding_bot/message"
	"github.com/bleemdy/ding_bot/schedule"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Bot struct {
	Messages       chan *Context
	Commands       map[string]func(ctx *Context)
	TextFnc        []func(ctx *Context)
	MessageManager *message.Manager
	Schedule       *schedule.Schedule
	WebHook        string
}

func (b *Bot) AddJob(taskName, spec string, f func(bot *Bot) func()) {
	b.Schedule.AddJob(taskName, spec, f(b))
}

func (b *Bot) push(msg *Context) {
	b.Messages <- msg
}

func (b *Bot) Send(content message.Common) {
	b.MessageManager.SendByHook(content)
}

func (b *Bot) handle() {
	list := [2]func(ctx *Context){b.commandHandle, b.textHandle}
	for msg := range b.Messages {
		for _, fn := range list {
			fn(msg)
		}
	}
}

func (b *Bot) OnCommand(reg string, f func(ctx *Context)) {
	b.Commands[reg] = f
}

func (b *Bot) commandHandle(c *Context) {
	for reg, fn := range b.Commands {
		matched, _ := regexp.MatchString(reg, c.Message.Text.Content)
		if matched {
			fn(c)
		}
	}
}

func (b *Bot) OnText(f func(ctx *Context)) {
	b.TextFnc = append(b.TextFnc, f)
}

func (b *Bot) textHandle(c *Context) {
	for _, fn := range b.TextFnc {
		fn(c)
	}
}

func (b *Bot) Run(addr, pattern string) {
	go b.MessageManager.Run()
	go b.Schedule.Run()
	go b.handle()
	http.HandleFunc(pattern, func(_ http.ResponseWriter, request *http.Request) {
		body, _ := ioutil.ReadAll(request.Body)
		msg := &message.Ding{}
		_ = json.Unmarshal(body, msg)
		ctx := newContext(b, msg)
		b.push(ctx)
	})
	fmt.Printf("runing in %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func New() *Bot {
	botManage := &Bot{
		Messages:       make(chan *Context, 10),
		Commands:       map[string]func(ctx *Context){},
		TextFnc:        []func(ctx *Context){},
		MessageManager: message.New(),
		Schedule:       schedule.New(),
	}
	return botManage
}
