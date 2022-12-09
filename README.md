# 自用钉钉机器人
1. 支持回调消息
2. 支持webhook主动发送
3. 定时任务（定时任务使用cron）

## 安装
```bash
    go get -u  github.com/bleemdy/ding_bot
```

## 基本使用
```go
    package main
    import (
        "github.com/bleemdy/ding_bot"
    )
    func main() {
        b := ding_bot.New()
        b.OnCommand("help", func(ctx *ding_bot.Context) {
            // do something...
        })
		b.OnText(func(ctx *ding_bot.Context) {
			// do something...
			// 会逐个执行执行添加的 text handler
		})
        b.Run("localhost:4000", "/ding")
    }
```

## 定时任务
```go
    package main
    import (
        "github.com/bleemdy/ding_bot"
        "github.com/bleemdy/ding_bot/message"
    )
    func main() {
        b := ding_bot.New()
        b.AddJob("task1", "@every 1s", func(bot *ding_bot.Bot) func() {
            return func() {
                bot.Send(
                    &message.Text{
                        Text:      "定时任务消息",
                        Webhook:   "webhook",
                    })
            }
		})
        b.Run("localhost:4000", "/ding")
    }
```