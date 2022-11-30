package context

import (
	"encoding/json"
	"github.com/bleemdy/ding_bot/message"
	"github.com/bleemdy/ding_bot/schedule"
	"io/ioutil"
	"log"
	"net/http"
)

type Bot struct {
	Messages       chan *Context
	Commands       map[string]func(ctx *Context)
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

func (b *Bot) commandHandle(cmd string, c *Context) {
	if f, ok := b.Commands[cmd]; ok {
		f(c)
	}
}

func (b *Bot) handle() {
	for msg := range b.Messages {
		l := len(msg.Args)
		if l > 0 {
			b.commandHandle(msg.Args[0], msg)
		}
	}
}

func (b *Bot) Command(c string, f func(ctx *Context)) {
	b.Commands[c] = f
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
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
