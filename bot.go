package ding_bot

import (
	"encoding/json"
	"fmt"
	"github.com/bleemdy/ding_bot/push"
	"github.com/bleemdy/ding_bot/schedule"
	"github.com/bleemdy/ding_bot/util"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type KeywordItem struct {
	pattern string
	fn      func(ctx *Context)
}

type Bot struct {
	Messages    chan *Context
	commandFns  map[string]func(*Context)
	messageFns  []func(*Context)
	KeywordFns  []*KeywordItem
	middleWares []func(*Context)
	pusher      push.Pusher
	scheduler   schedule.Scheduler
}

func (b *Bot) Use(f ...func(*Context)) {
	b.middleWares = append(b.middleWares, f...)
}

func (b *Bot) Send(content push.Common) {
	b.pusher.Send(content)
}

func (b *Bot) SendByHook(content push.Common) {
	b.pusher.SendByHook(content)
}

func (b *Bot) SetPush(messageManager push.Pusher) {
	b.pusher = messageManager
}

func (b *Bot) SetScheduler(schedule schedule.Scheduler) {
	b.scheduler = schedule
}

func (b *Bot) AddJob(taskName, spec string, f func(bot *Bot) func()) {
	b.scheduler.AddJob(taskName, spec, f(b))
}

func (b *Bot) OnCommand(reg string, f func(ctx *Context)) {
	b.commandFns[reg] = f
}

func (b *Bot) onCommandsHandle(c *Context) {
	if c.Command == "" {
		return
	}
	if fn, ok := b.commandFns[c.Command]; ok {
		fn(c)
	}
}

func (b *Bot) OnKeyword(pattern string, fs ...func(ctx *Context)) {
	for _, fn := range fs {
		b.KeywordFns = append(b.KeywordFns, &KeywordItem{
			pattern: pattern,
			fn:      fn,
		})
	}
}

func (b *Bot) onKeywordsHandle(c *Context) {
	for _, item := range b.KeywordFns {
		matched, _ := regexp.MatchString(item.pattern, c.Content)
		if matched {
			item.fn(c)
		}
	}
}

func (b *Bot) OnMessage(f ...func(ctx *Context)) {
	b.messageFns = append(b.messageFns, f...)
}

func (b *Bot) onMessagesHandle(c *Context) {
	for _, fn := range b.messageFns {
		fn(c)
	}
}

func (b *Bot) push(msg *Context) {
	b.Messages <- msg
}

func (b *Bot) handle() {
	list := [3]func(ctx *Context){b.onCommandsHandle, b.onMessagesHandle, b.onKeywordsHandle}
	for msg := range b.Messages {
		msg.Run()
		if msg.IsAborted() {
			continue
		}
		for _, fn := range list {
			fn(msg)
		}
	}
}

func (b *Bot) Run(addr, pattern string) {
	go b.pusher.Run()
	go b.scheduler.Run()
	go b.handle()
	http.HandleFunc(pattern, func(ResponseWriter http.ResponseWriter, request *http.Request) {
		body, _ := ioutil.ReadAll(request.Body)
		msg := &push.Ding{}
		_ = json.Unmarshal(body, msg)
		ctx := newContext(msg)
		if ctx == nil {
			return
		}
		ctx.handlers = b.middleWares
		ctx.Bot = b
		ctx.ResponseWriter = ResponseWriter
		ctx.Request = request
		b.push(ctx)
	})
	fmt.Printf("runing in %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func New() *Bot {
	botManage := &Bot{
		Messages:    make(chan *Context, 10),
		commandFns:  map[string]func(ctx *Context){},
		messageFns:  []func(ctx *Context){},
		middleWares: []func(*Context){},
		pusher:      push.New(),
		scheduler:   schedule.New(),
		KeywordFns:  []*KeywordItem{},
	}
	return botManage
}
func NewDefault() *Bot {
	util.ParseConfig()
	if viper.GetString("appSecret") == "" {
		log.Fatal("请配置config.toml appSecret")
	}
	botManage := &Bot{
		Messages:    make(chan *Context, 10),
		commandFns:  map[string]func(ctx *Context){},
		messageFns:  []func(ctx *Context){},
		middleWares: []func(*Context){},
		pusher:      push.New(),
		scheduler:   schedule.New(),
		KeywordFns:  []*KeywordItem{},
	}
	botManage.Use(VerifyRequest)
	return botManage
}
