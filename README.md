# 自用钉钉机器人
1. 支持回调消息
2. 支持webhook主动发送
3. 定时任务（定时任务使用[go cron](github.com/robfig/cron)）
4. 中间件

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
        b.OnCommand("/help", func(ctx *ding_bot.Context) {
            // do something...
            // 命令以 / 开头
        })
        b.OnMessage(func(ctx *ding_bot.Context) {
            // do something...
            // 会逐个执行执行添加的 message handler
        }) 
        b.OnKeyword("test",func(ctx *ding_bot.Context) {
            // do something...
            // 会逐个执行所有匹配 keyword handler
        })
        b.Run("localhost:4000", "/ding")
    }
```

## 定时任务

```go
    package main
    
    import (
        "github.com/bleemdy/ding_bot"
        "github.com/bleemdy/ding_bot/push"
    )

    func main() {
        b := ding_bot.New()
        b.AddJob("task1", "@every 1s", func(bot *ding_bot.Bot) func() {
            return func() {
                bot.Send(
                    &push.Text{
                        Text:    "定时任务消息",
                        Webhook: "webhook",
                    })
            }
        })
        b.Run("localhost:4000", "/ding")
    }
```
## 中间件

```go
    package main
    import (
        "github.com/bleemdy/ding_bot"
    )
    func main() {
        b := ding_bot.New()
        b.Use(func(ctx *ding_bot.Context) {
			//ctx.Next()
			//or
			//ctx.Abort()
        })
        b.Run("localhost:4000", "/ding")
    }
```


## 带请求合法性验证
```go
    package main
    import (
        "github.com/bleemdy/ding_bot"
    )
    b := ding_bot.NewDefault()
    // 必须在config.toml中配置机器人的appSecret
    b.Run("localhost:4000", "/ding")
```