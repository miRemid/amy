package server

import (
	"fmt"
	"io"
	"log"
	"context"
	"io/ioutil"
	"net/http"
	"encoding/json"
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
		// 将上下文打包成一个event
		bytedata, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		// 封装事件
		ctx := context.WithValue(r.Context(), WriterKey, &w)
		ctx = context.WithValue(ctx, RequestKey, r)
		ctx = context.WithValue(ctx, ByteData, bytedata)
		event := newEvent(ctx, bot)
		event.handler = middle
		event.hlength = len(middle)
		// 传递中间件
		event.Next()
	}
}

// ParseMessage 解析消息转发到不同的Handler中
func (bot *Bot) ParseMessage(event *CQEvent){
	bytedata := event.Value(ByteData).([]byte)
	var msg map[string]interface{}
	json.Unmarshal(bytedata, &msg)
	event.Type = msg["post_type"].(string)
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
	}
	return
}

// SignatureMiddleware CQHTTP消息验证中间
func (bot Bot) signature(event *CQEvent){
	sig := event.reqHeader().Get("X-Signature")
	if sig == "" {
		log.Println("未找到头部信息")
		return
	}
	sig = sig[len("sha1="):]
	mac := hmac.New(sha1.New, []byte(bot.scret))
	byteData, _ := event.Value(ByteData).([]byte)

	io.WriteString(mac, string(byteData))
	res := fmt.Sprintf("%x", mac.Sum(nil))
	log.Printf("CQ HMAC:%s, Amy HMAC:%s\n", sig, res)
	if res == sig {
		log.Println("消息来自酷Q")
	}else{
		log.Println("消息不来自酷Q，已屏蔽处理")
		return
	}		
	event.Next()
}