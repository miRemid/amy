# Amy
Amy是一个轻量级cqhttp的go版sdk，目前使用文档较乱，将会逐步改正
# install
使用前请安装[酷Q](https://cqp.cc/)和[CQHTTP](https://cqhttp.cc/docs/4.11/#/)
```
git clone https://github.com/miRemid/amy.git
```
# Usage
## Send Private Message
```golang
import "github.com/miRemid/amy"
import cqmsg "github.com/miRemid/amy/message"
import "time"
func main(){
    // 创建api
    api := amy.NewAmyAPI("localhost", 5700)
    // 创建消息生成器
    builder := message.NewCQMsgBuilder()
    // 创建一条私人信息
    msg := builder.PrivateMsg(123456789, "test", true)
    // 发送私人信息
	if res, ok := api.SendPrivateMsg(msg, false); ok {
		fmt.Println(res.ID)
	}else{
		fmt.Println("Send Failed")
    }
    // 检查能否发送图片
	if ok := api.CanSendImage(false); ok {
		fmt.Println("Can")
	}else{
		fmt.Println("No")
	}
}
```
其中，消息不仅限于字符串形式，推荐字符串形式发送
```
msg := "String"
msg := cqmsg.CQJSON{
    "type":"text",
    "data":{"text":"test"}
}
// 数组模式正在完善中
msg := cqmsg.CQJSON{
    cqmsg.CQJSON{
            "type":"text",
            "data":JSONMsg{
            "text":"test1",
        }
    },
    cqmsg.CQJSON{
            "type":"text",
            "data":JSONMsg{
            "text":"test2",
        }
    },
}
```
具体配置详见[CQHTTP文档](https://cqhttp.cc/docs/4.11/#/Message)
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
            // 可以快速回复，响应数据详见CQHTTP文档
            event.JSON(message.CQJSON{
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
# TODO
- 覆盖CQHTTP所有常用HTTP API(已完成)
- 创建轻量级Serve端(v0.0.1)
- AmyMQ消息队列处理服务
