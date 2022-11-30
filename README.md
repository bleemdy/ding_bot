# 自用钉钉机器人
1. 支持回调消息
2. 支持webhook主动发送
3. 定时任务（定时任务使用cron）

## 基本使用
```go
    package main
    import (
        "github.com/bleemdy/ding_bot/context"
    )
    func main() {
        b := context.New()
        b.Command("help", func(ctx *context.Context) {
            // do something...
        })
        b.Run("localhost:4000", "/ding")
    }
```

## 定时任务
```go
    package main
    import (
        "github.com/bleemdy/ding_bot/context"
		"github.com/bleemdy/ding_bot/message"
    )
    func main() {
        b := context.New()
        b.SetWebHook("webhook url")
        b.AddJob("task1", "@every 1s", func(bot *context.Bot) func() {
            return func() {
				bot.Send(
                    &message.Text{
                        Basic: message.Basic{
                            Text:      "定时任务消息",
                            Webhook:   b.WebHook(),
                        },
                    })
            }
		})
        b.Run("localhost:4000", "/ding")
    }
```