package push

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type Pusher interface {
	Send(Common)
	SendByHook(Common)
	Run()
}

type Push struct {
	message       chan Common
	activeMessage chan Common
}

func (m Push) Send(msg Common) {
	m.message <- msg
}
func (m Push) SendByHook(msg Common) {
	m.activeMessage <- msg
}

func (m Push) Run() {
	go func() {
		for mes := range m.message {
			go m.sendFunc(mes)
		}
	}()
	go func() {
		for mes := range m.activeMessage {
			go m.sendFunc(mes)
			time.Sleep(3 * time.Second)
		}
	}()
}

func (m Push) sendFunc(msg Common) {
	url, body, header := msg.getContent()
	client := resty.New().R()
	client.SetHeaders(header)
	client.SetBody(body)
	result := &Result{}
	res, _ := client.SetResult(result).Post(url)
	if result.Errmsg != "ok" {
		fmt.Println(res)
	}
}

func New() *Push {
	return &Push{
		message:       make(chan Common, 6),
		activeMessage: make(chan Common, 6),
	}
}
