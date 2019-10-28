package model

import (
	"github.com/miRemid/amy/utils"
)

// CQEvent is the event from cqhttp's websocket server
type CQEvent struct {
	Body 	[]byte
	Map	 	map[string]interface{}
	Type 	string
}

// NewCQEvent packs the byte data into a CQEvent
func NewCQEvent(body []byte) CQEvent{
	msg := utils.LoadIntoMap(body)
	return CQEvent{
		Body: 	body,
		Type: 	msg["post_type"].(string),
		Map:	msg,
	}
}	