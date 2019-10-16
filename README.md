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
# TODO
- 覆盖CQHTTP所有常用HTTP API(已完成)
- 创建轻量级Serve端(v0.0.1)
- AmyMQ消息队列处理服务
