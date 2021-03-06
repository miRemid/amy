Amy是一个轻量级cqhttp的go版sdk，目前使用文档较乱，将会逐步整理
- [安装](#安装)
- [使用](#使用)
    - [初步](#初步)
    - [消息格式](#消息格式)
    - [CQ码](#cq码)
- [服务端](#服务端)
    - [Http:v0.0.1](#http:v0.0.1)
    - [Http:v0.0.2](#http:v0.0.2)
    - [WebSocket](#websocket)
    - [AmyMQ](#amymq)
- [TODO](#todo)
# 安装
使用前请安装[酷Q](https://cqp.cc/)和[CQHTTP](https://cqhttp.cc/docs/4.11/#/)
```
git clone https://github.com/miRemid/amy.git
```
# 使用
## 初步
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
    // 发送给123456789，消息为”test"，true表明消息是否需要转义，仅字符串消息有效
    msg := builder.PrivateMsg(123456789, "test", true)
    // 发送私人信息
    if res, ok := api.SendPrivateMsg(msg, false); ok {
        fmt.Println(res.ID)
    }else{
        fmt.Println("Send Failed")
    }

    // 当然也可以直接通过flag发送消息
    // 1234565为目标id，"test"为消息内容，true为消息是否需要转义仅字符串有效，false表示是否调用异步api，amy.Private标识该消息类型为私人消息
    api.Send(1234565, "test", true, false, amy.Private)
    // 检查能否发送图片
    if ok := api.CanSendImage(false); ok {
        fmt.Println("Can")
    }else{
        fmt.Println("No")
    }
}
```
## 消息格式
消息不仅限于字符串形式
```golang
// 字符串
msg := "String"
// 消息块
/*
{
    "type":"text",
    "data":{
        "text":"hello"
    }
}
*/
msg := cqmsg.CQJSON("text", "text", "hello")
// 消息块数组
/*
[
    {
        "type":"text",
        "data":{
            "text":"hello"
        }
    },{
        "type":"face",
        "data":{
            "id":"111"
        }
    }
]
*/
msg := cqmsg.CQArray{
    cqmsg.CQJSON("text", "text", "hello"),
    cqmsg.CQJSON("face", "id", "111"),
}
```
具体配置详见[CQHTTP文档](https://cqhttp.cc/docs/4.11/#/Message)
## CQ码
Amy中可以生成CQ码
```golang
import "github.com/miRemid/amy/cqcode"

// [CQ:text,file=asdf]
cq := cqcode.CqCode("text", cqcode.CQParams{
    "file": "asdf",
})
// [CQ:face,id=1]
face := cqcode.Face(1)
```
Amy中提供了解析CQ码的函数`cqcode.CQSplit`
```golang
msg := "asdf[CQ:at,id=fdfsdf]asdf[CQ:at,id=asdfdf]dasfioghfdguihn[CQ:at,id=ghjgfhj]"
// cqs -> []cqcode.CQCdoee
cqs := cqcode.CQSplit(msg)
for _, cqcode := range cqs {
    fmt.Printf("func=%v,params=%v\n", cqcode.Func, cqcode.Params)
}
```
# 服务端
## Http:v0.0.1
在`amy/server`中可以创建一个小型服务器，具体请见[server](https://github.com/miRemid/amy/tree/master/server)
## Http:v0.0.2
在`amy/tserver`中可以建立一个小型服务器，具体请见[tserver](https://github.com/miRemid/amy/tree/master/tserver)
这里有一个现成的利用`tserver`建立的Amy框架[Hanabi](https://github.com/miRemid/amy/hanabi)
## WebSocket
已支持websocket，`github.com/miRemid/amy/websocket`
```golang
import "github.com/miRemid/amy/websocket"
import "github.com/miRemid/amy/websocket/model"
import "log"
func main(){
    // 创建总客户端，默认api和event地址端口一致，token为空
    client := websocket.NewCQClient("127.0.0.1", 6700)
    // 设置api响应处理函数
    client.OnResponse(func(res model.CQResponse){
        log.Println(res.Data)
    })
    // 单独设置api地址端口
    // client.SetAPIConfig(apiurl, apiport)
    // 设置api的token
    // client.SetToken("token")
    // 设置event响应处理函数
    client.OnMessage(func(evt model.CQEvent){
        if msg := evt.Map["raw_message"].(string); msg == "hello" {
            if t := evt.Map["message_type"].(string); t == "private"{
                client.Send("send_private_msg", model.CQParams{
                    "user_id": 123456,
                    "message": "hello",
                })
            }
        }
    })
    client.Run()
}
```
你也可以单独创建消息接受和API发送websocket client
```golang
import "github.com/miRemid/amy/websocket"
import "github.com/miRemid/amy/websocket/model"
import "log"

func main(){
    // api客户端
    api := websocket.NewAPIClient("127.0.0.1", 6700, "")
    api.OnResponse(func(evt model.CQResponse){
        log.Printf(evt.Status)
    })
    // 消息接收客户端
    client := websocket.NewClient("127.0.0.1", 6700)
    client.OnMessage(func(evt model.CQEvent){                
        if msg := evt.Map["raw_message"].(string); msg == "hello" {
            if t := evt.Map["message_type"].(string); t == "private"{
                // 每次发送会创建一次连接
                go api.Send("send_private_msg", model.CQParams{
                    "user_id": 123456,
                    "message": "hello",
                })
                api.Send("send_private_msg", model.CQParams{
                    "user_id": 123456,
                    "message": "hello",
                })
            }
        }
    })
    client.Run()
}

```
# AmyMQ
可以在Release中下载AmyMQ进行消息队列转发，请按照`amy/amymq`文件夹中的config进行配置.
AmyMQ目前还在完善中，只适配英文开头的标准命令格式`[!,#,...]cmd params`，消息转发过程如下：
```
发送消息: !hello 你好
CQHTTP: 接受消息->转发到AmyMQ
AmyMQ: 接受消息->处理消息加入队列->分析消息(cmd:hello)->转发消息到"http://你的AmyMQ配置项/hello"
```
[config.json配置项详情](https://github.com/miRemid/amy/tree/master/amymq)
# TODO
- 覆盖CQHTTP所有常用HTTP API(已完成)
- 创建Serve端(v0.0.2)
- AmyMQ消息队列处理服务(v0.0.1)
