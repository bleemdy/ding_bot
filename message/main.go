package message

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type Manager struct {
	message       chan Common
	activeMessage chan Common
}

func (m *Manager) Send(msg Common) {
	m.message <- msg
}
func (m *Manager) SendByHook(msg Common) {
	m.activeMessage <- msg
}

func (m Manager) Run() {
	go func() {
		for mes := range m.message {
			go m.send(mes)
		}
	}()
	go func() {
		for mes := range m.activeMessage {
			m.send(mes)
			time.Sleep(3 * time.Second)
		}
	}()
}

func (m Manager) send(msg Common) {
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

func New() *Manager {
	return &Manager{
		message:       make(chan Common, 6),
		activeMessage: make(chan Common, 6),
	}
}
