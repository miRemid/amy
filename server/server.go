package server

// 轻量级CQHTTP后端服务

import (
	"net/http"
	"log"
	"fmt"

	"github.com/miRemid/amy/message"

	"github.com/gorilla/mux"
)

const (
	// KMessage 普通消息
	KMessage = "message"
	// KNotice 通知消息
	KNotice = "notice"
	// KRequest 请求消息
	KRequest = "request"
)

// EventMap map
type EventMap map[string]interface{}

// Bot 机器人对象
type Bot struct {
	scret string					// CQHTTP配置项
	router *mux.Router				// 路由
	handlers []CQEventHandler		// 中间件
	parse	CQEventHandler			// 最终消息处理函数
	messageHandler 	CQEventHandler	// 普通消息处理函数
	noticeHandler 	CQEventHandler	// 提示消息处理函数
	requestHandler 	CQEventHandler	// 请求消息处理函数
}

// Hello 你好
func Hello(event CQEvent) {
	msg := event.Map()
	if msg["raw_message"] == "你好" {
		event.JSON(200, message.CQJSON{
			"reply":"You too~~",
		})
	}	
}

// NewServer 实例化一个Bot对象
// addr: cqhttp api域名
// port: cqhttp api端口
func NewServer() *Bot {
	var res Bot
	res.router = mux.NewRouter()
	res.messageHandler = Hello
	res.parse = res.ParseMessage
	res.noticeHandler = func(event CQEvent){event.JSON(204, nil)}
	res.requestHandler = func(event CQEvent){event.JSON(204, nil)}
	return &res
}

// Use 构造中间件链
func (bot *Bot) Use(handlers ...CQEventHandler) {
	bot.handlers = append(bot.handlers, handlers...)	
}

// Signature 开启验证
func (bot *Bot) Signature(key string){
	log.Printf("已开启消息来源验证, Screct=%s\n", key)
	bot.scret = key
	bot.Use(bot.signature)
}

// SetParse 设置消息最终处理函数
func (bot *Bot) SetParse(handler CQEventHandler) {
	bot.parse = handler
}

// Register 注册函数
func (bot *Bot) Register(name string, handler CQEventHandler) error{
	switch name {
	case KMessage:
		bot.messageHandler = handler
		break
	case KNotice:
		bot.noticeHandler = handler
		break
	case KRequest:
		bot.requestHandler = handler
		break
	default:
		return fmt.Errorf("%s not found", name)
	}
	return nil
}

// Run 建立一个http服务
// 在开启服务之前，需要绑定事件上报处理函数
// addr: 运行域名
// port: 运行端口
func (bot *Bot) Run(addr string, router string) {
	if len(addr) <= 7{
		log.Printf("Amy Server Run At http://0.0.0.0%s\n", addr)
	}else{
		log.Printf("Amy Server Run At http://%s\n", addr)		
	}
	bot.Use(bot.parse)
	bot.router.HandleFunc(router, bot.pack(bot.handlers)).Methods("POST", "GET")
	err := http.ListenAndServe(addr, bot.router)
	if err != nil {
		log.Fatal("Failed to start HTTP server")
		return
	}
}