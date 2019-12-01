package tserver

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/miRemid/amy"
	"github.com/miRemid/amy/message"
	"github.com/miRemid/amy/tserver/event"
	"github.com/miRemid/amy/utils"
)

func (bot *Bot) pack(handlers ...event.CQEventHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		base := &event.CQEventBase{
			API: bot.apipool.Get().(*amy.API),
		}
		base.SetHTTP(w, r)
		data, _ := ioutil.ReadAll(r.Body)
		base.Body = data
		tmsg := utils.LoadIntoMap(data)
		// 判断消息类型
		postType, ok := tmsg["post_type"].(string)
		if !ok {
			log.Println("post_type查找失败")
		}
		switch postType {
		case "message":
			session := toCQSession(tmsg)
			session.CQEventBase = base
			session.Use(handlers...)
			log.Println("CQParse=", session.CQCode())
			session.Next()
			break
		case "notice":
			break
		case "request":
			break
		}
	}
}

func toCQSession(tmsg map[string]interface{}) *event.CQSession {
	var res event.CQSession

	tsender, ok := tmsg["sender"].(map[string]interface{})
	if !ok {
		return &res
	}
	var sender message.CQSender
	sender.NickName = tsender["nickname"].(string)
	sender.Sex = tsender["sex"].(string)
	userid, _ := tsender["user_id"].(json.Number).Int64()
	age, _ := tsender["age"].(json.Number).Int64()
	sender.Age = int32(age)
	sender.UserID = int(userid)
	res.Sender = sender

	res.Type = tmsg["message_type"].(string)
	res.Message = tmsg["message"].(string)
	res.RawMessage = tmsg["raw_message"].(string)

	return &res
}

func Signature(key string) event.CQEventHandler {
	log.Printf("以开启Signature验证, key=%v\n", key)
	return func(evt event.CQEvent) {
		sig := evt.ReqHeader().Get("X-Signature")
		if sig == "" {
			log.Println("未找到X-Signature头部信息，请检查CQHTTP配置")
			evt.JSON(204, nil)
			return
		}
		sig = sig[len("sha1="):]
		mac := hmac.New(sha1.New, []byte(key))
		byteData := evt.GetBody()

		io.WriteString(mac, string(byteData))
		res := fmt.Sprintf("%x", mac.Sum(nil))
		if res != sig {
			log.Println("消息不来自CQHTTP，以屏蔽处理")
			evt.JSON(204, nil)
		} else {
			log.Println("接受到CQHTTP消息，开始解析处理")
			evt.Next()
		}
	}
}
