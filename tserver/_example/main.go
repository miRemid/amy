package main

import (
	"log"

	"github.com/miRemid/amy/tserver"
	"github.com/miRemid/amy/tserver/event"

)

func test(evt event.CQSession) {
	evt.Send("hello", true, true)
}

func main() {
	bot := tserver.NewBot("127.0.0.1", 5700)
	bot.AccessToken = "asdf"
	bot.Use(tserver.Signature("amy"))
	bot.On(test, tserver.Message)
	log.Println("listen localhost:3000")
	bot.Run(":3000", "/")
}
