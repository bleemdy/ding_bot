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
	messages       chan *Context
	commands       map[string]func(ctx *Context)
	messageManager *message.Manager
	schedule       *schedule.Schedule
	webHook        string
}

func (b *Bot) WebHook() string {
	return b.webHook
}

func (b *Bot) AddJob(taskName, spec string, f func(bot *Bot) func()) {
	b.schedule.AddJob(taskName, spec, f(b))
}

func (b *Bot) push(msg *Context) {
	b.messages <- msg
}

func (b *Bot) Send(content message.Common) {
	b.messageManager.SendByHook(content)
}

func (b *Bot) commandHandle(cmd string, c *Context) {
	if f, ok := b.commands[cmd]; ok {
		f(c)
	}
}

func (b *Bot) handle() {
	for msg := range b.messages {
		l := len(msg.Args)
		if l > 0 {
			b.commandHandle(msg.Args[0], msg)
		}
	}
}

func (b *Bot) Command(c string, f func(ctx *Context)) {
	b.commands[c] = f
}

func (b *Bot) SetWebHook(webHook string) {
	b.webHook = webHook
}

func (b *Bot) Run(addr, pattern string) {
	go b.messageManager.Run()
	go b.schedule.Run()
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

func New() *Bot {

	botManage := &Bot{
		make(chan *Context, 10),
		map[string]func(ctx *Context){},
		message.New(),
		schedule.New(),
		"",
	}
	return botManage
}
