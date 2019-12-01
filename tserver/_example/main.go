package main

import (
	"log"

	"github.com/miRemid/amy/tserver"
)

func main() {
	bot := tserver.NewBot("127.0.0.1", 5700)
	bot.Screct = "amy"
	log.Println("listen localhost:3000")
	bot.Run(":3000", "/")
}
