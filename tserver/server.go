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

type Bot struct {
	Screct      string
	AccessToken string
	Timeout     int

	handlers []event.CQEventHandler
	router   *mux.Router

	apipool *sync.Pool
	apiURL  string
	apiport int

	MessageHandler event.CQEventHandler
	NoticeHandler  event.CQEventHandler
	RequestHandler event.CQEventHandler
}

func NewBot(apiURL string, port int) *Bot {
	var res Bot
	res.handlers = make([]event.CQEventHandler, 0)
	res.router = mux.NewRouter()
	res.apiURL = apiURL
	res.apiport = port
	return &res
}

func (bot *Bot) Use(handler event.CQEventHandler) {
	bot.handlers = append(bot.handlers, handler)
}

func (bot *Bot) Run(addr string, router string) {
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
	if bot.Screct != "" {
		bot.Use(Signature(bot.Screct))
	}
	bot.router.HandleFunc(router, bot.pack(bot.handlers...)).Methods("POST", "GET")
	err := http.ListenAndServe(addr, bot.router)
	if err != nil {
		log.Fatal("Failed to start HTTP server")
		return
	}
}
