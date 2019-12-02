package tserver

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/miRemid/amy"

	"github.com/gorilla/mux"
	"github.com/miRemid/amy/tserver/event"
)

const (
	// Message Event Flag
	Message = iota
	// Notice Event Flag
	Notice
	// Request Event Flag
	Request
)

// Bot client
type Bot struct {
	AccessToken string
	Timeout     int

	handlers []event.CQEventHandler
	router   *mux.Router

	apipool *sync.Pool
	apiURL  string
	apiport int

	messageHandler CQSessionHandler
	noticeHandler  CQNoticeHandler
	requestHandler CQRequestHandler
}

// NewBot return a Bot client
func NewBot(apiURL string, port int) *Bot {
	var res Bot
	res.handlers = make([]event.CQEventHandler, 0)
	res.router = mux.NewRouter()
	res.apiURL = apiURL
	res.apiport = port

	res.messageHandler = messageHandler
	res.noticeHandler = noticeHandler
	res.requestHandler = requestHandler

	return &res
}

// Use sevral handlers as middlware
func (bot *Bot) Use(handler ...event.CQEventHandler) {
	bot.handlers = append(bot.handlers, handler...)
}

// On set the event hanlder, flag should from tserver.Message,tserver.Notice,tserver.Request
func (bot *Bot) On(handler interface{}, flag int) {
	switch flag {
	case Message:
		if middle, ok := handler.(func (evt event.CQSession)); ok {
			bot.messageHandler = middle
		}else {
			log.Fatal("Handler参数错误")
		}
		break
	case Notice:
		if middle, ok := handler.(func (evt event.CQNotice)); ok {
			bot.noticeHandler = middle
		}else {
			log.Fatal("Handler参数错误")
		}
		break
	case Request:
		if middle, ok := handler.(func (evt event.CQRequest)); ok {
			bot.requestHandler = middle
		}else {
			log.Fatal("Handler参数错误")
		}
		break
	}
}

// Run a http server at addr/router
func (bot *Bot) Run(addr string, router string, handlers ...event.CQEventHandler) {
	bot.Use(handlers...)
	bot.apipool = &sync.Pool{
		New: func() interface{} {
			res := amy.NewAmyAPI(bot.apiURL, bot.apiport)
			if bot.AccessToken != "" {
				res.SetToken(bot.AccessToken)
			}
			res.SetTimeout(time.Second * time.Duration(bot.Timeout))
			return res
		},
	}
	bot.Use(bot.convert)
	bot.router.HandleFunc(router, bot.pack(bot.handlers...)).Methods("POST", "GET")
	err := http.ListenAndServe(addr, bot.router)
	if err != nil {
		log.Fatal("Failed to start HTTP server")
		return
	}
}
