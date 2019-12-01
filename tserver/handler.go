package tserver

import (
	"github.com/miRemid/amy/tserver/event"
)

func messageHandler(evt event.CQSession) {
	evt.JSON(200, event.Map{
		"reply":       evt.RawMessage,
		"auto_escape": true,
	})
}

func noticeHandler(evt event.CQNotice) {
	evt.JSON(200, nil)
}

func requestHandler(evt event.CQRequest) {
	evt.JSON(200, nil)
}
