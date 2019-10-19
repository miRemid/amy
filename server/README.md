# Server服务端demo
## 普通
在server包中，自带了一个轻量级服务端，你可以创建一个小型处理服务器，服务器默认使用`server.parse`函数作为最终处理函数
```golang
import (
    "amy/message"
    "amy/server"
)
// MessageHandler 普通消息处理
func MessageHandler(event server.CQEvent) {
    // 获取消息类型，仅在post_type为message有效
    msgtype := event.MessageType()
    if msgtype == "private"{
        // 创建一条私人消息
        var msg message.CQPrivate
        event.ReadJSON(&msg)
        // 如果消息为”你好“
        if msg.RawMessage == "你好"{
            // 响应CQHTTP
            event.JSON(200, message.CQMAP{
                "reply": "You too~~",
            })
        }
    }
}

func main() {
    // 创建一个服务，参数为酷Q的api地址
    bot := server.NewServer("127.0.0.1", 5700)
    // 注册普通消息处理函数
    bot.Register(server.KMessage, MessageHandler)
    // 运行在localhost:3000，酷Q上报路由为"/"
    bot.Run(":3000", "/")
}
```
### 注意
如果用户在处理业务逻辑时没有响应CQHTTP，且服务端的Parse函数为默认值，会默认返回204给CQHTTP。`bot.Register`注册消息处理的函数只有在执行`bot.SetParse(bot.ParseMessage)`或默认时有效。
```golang
// SetParse 设置消息最初处理函数
// 自带的ParseMessage会将消息进行分层，转发到普通消息、提示、请求处理函数中
// 通过bot.Register注册以上三个函数
bot.SetParse(bot.ParseMessage)
```
## 中间件
你可以对每一个事件注册中间件处理，比如消息的认证，判定权限等，需要注意的是，注册中间件需要按照顺序执行，切必须在注册函数前执行,中间件结构为`func Handler(event server.CQEvent)`
```golang
import (
    "amy/message"
    "amy/server"
)
func main() {
    // 创建一个服务，参数为酷Q的api地址
    bot := server.NewServer("127.0.0.1", 5700)
    // 添加CQ消息判断中间件    
    bot.Signature("amy")
    bot.Use(Your CQEvent Handler)
    // 注册普通消息处理函数
    bot.Register(server.KMessgae, server.Hello)
    // 运行在localhost:3000，酷Q上报路由为"/"
    bot.Run(":3000", "/")
}
```