package server

import (
	"fmt"
	"io"
	"log"
	"context"
	"io/ioutil"
	"net/http"
	"crypto/sha1"
	"crypto/hmac"	
)

type key string

const (
	// ByteData key of []byte value
	ByteData 	key = "byte"
	// EventKey key of event
	EventKey 	key = "event"
	// WriterKey key of writer
	WriterKey 	key = "writer"
	// RequestKey key of request
	RequestKey 	key = "request"
)

func (bot Bot) pack(middle []CQEventHandler) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		bytedata, _ := ioutil.ReadAll(r.Body)
		msg := loadintomap(bytedata)
		defer r.Body.Close()
		ctx := context.WithValue(r.Context(), ByteData, bytedata)
		event := CQEvent{
			Type: msg["post_type"].(string),
			ctx: ctx,
			writer: w,
			req: r,
			handler: middle,
			hlength: len(middle),
		}		
		event.Next()
	}
}

// ParseMessage 解析消息转发到不同的Handler中
func (bot *Bot) ParseMessage(event CQEvent){
	switch event.Type {
	case "message":
		bot.messageHandler(event)
		break
	case "notice":
		bot.noticeHandler(event)
		break
	case "request":
		bot.requestHandler(event)
		break
	default:
		event.JSON(204, nil)
		break
	}
	return
}

// SignatureMiddleware CQHTTP消息验证中间件
func (bot Bot) signature(event CQEvent){
	sig := event.reqHeader().Get("X-Signature")
	if sig == "" {
		log.Println("未找到头部信息，请检查CQHTTP配置")
		event.JSON(204, nil)
		return
	}
	sig = sig[len("sha1="):]
	mac := hmac.New(sha1.New, []byte(bot.scret))
	byteData, _ := event.Value(ByteData).([]byte)

	io.WriteString(mac, string(byteData))
	res := fmt.Sprintf("%x", mac.Sum(nil))
	if res != sig {
		log.Println("消息不来自酷Q，已屏蔽处理")
		event.JSON(204, nil)
		return
	}
	event.Next()
}