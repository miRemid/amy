package tserver

import (
	"github.com/miRemid/amy/tserver/event"
)

// CQSessionHandler is the function solve the CQHTTP message event
type CQSessionHandler func (evt event.CQSession)
// CQNoticeHandler is the function solve the CQHTTP notice event
type CQNoticeHandler func (evt event.CQNotice)
// CQRequestHandler is the function solve the CQHTTP request event
type CQRequestHandler func (evt event.CQRequest)

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
