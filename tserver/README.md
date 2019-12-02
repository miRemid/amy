# HTTP:v0.0.2(tserver)

# 快速使用
```golang
package main

import (
	"log"

	"github.com/miRemid/amy/tserver"
	"github.com/miRemid/amy/tserver/event"

)

// 消息处理函数
func test(evt event.CQSession) {
    evt.Send("hello", true, true)
}

func main() {
    // 创建Bot
    bot := tserver.NewBot("127.0.0.1", 5700)
    // 设置Amy Sdk Token，与cqhttp配置一致，如果没有则无需设置
    bot.AccessToken = "asdf"
    // 使用Signature中间件，中间件参数会在下面阐述
    bot.Use(tserver.Signature("amy"))
    // 设置CQHTTP Event解析函数
    // tserver.Message为函数标志符，可选如下
    // tserver.Message; tserver.Notice; tserver.Request
    // 每个解析函数的参数是不一致的
    // Message的handler为func(evt event.CQSession)
    // Notice的handler为func(evt event.CQNotice)
    // Request的handler为func(evt event.CQRequest)
    bot.On(test, tserver.Message)
    log.Println("listen localhost:3000")
    // 监听本地3000端口，路由为"/"
    bot.Run(":3000", "/")
}
```
# 中间件
你可以使用`bot.Use`来添加消息传递的中间件，中间件函数模板为`func middleware(evt event.CQEvent)`与Message,Request,Notice的解析函数不一致，请格外注意。
并且中间件使用的顺序是顺序处理的，基于先来先到原则，在使用过程中务必注意顺序