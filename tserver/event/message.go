package event

import (
	"encoding/json"
	"regexp"

	"github.com/miRemid/amy/cqcode"

	"github.com/miRemid/amy"

	"github.com/miRemid/amy/message"
	"github.com/miRemid/amy/tserver/utils"
)

var (
	reg = regexp.MustCompile(`\[CQ:(.*?)\]`)
)

// CQSession is the message event struct
type CQSession struct {
	// CQEvent is the base event
	*CQEvent
	// Type is the message's type
	// Private;Group;Discuss
	Type       string
	// Sender is the message sender's infomation struct
	Sender     message.CQSender
	// Message is the cqhttp parse message from raw qq message
	Message    string
	// RawMessage is the raw qq message
	RawMessage string
}

// Params will try to parse event's message by cmds
// It will return two string, cmd and params
// Example: !help Test, give cmds as ["!"], will return help, Test
func (evt CQSession) Params(cmds ...string) (cmd, params string) {
	msg := string(evt.Body)
	return utils.CmdParser(msg, cmds...)
}

// CQCode will try to parse event's message
// It will return a string array containers CQCode like [CQ:face,id=1]
func (evt CQSession) CQCode() (res []cqcode.CQCode) {
	res = make([]cqcode.CQCode, 0)
	str := reg.FindAllString(evt.Message, -1)
	for _, v := range str {
		res = append(res, cqcode.CQParse(v))
	}
	return
}

// ReadJSON will loads event's []byte data into a struct like message.CQPrivate or message.CQGroup
func (evt CQSession) ReadJSON(cq interface{}) error {
	body := evt.Body
	return json.Unmarshal(body, cq)
}

// Send can send a msg to CQHTTP server by event.Type
func (evt CQSession) Send(msg interface{}, auto, async bool) (res message.CQMessageID, err error) {
	switch evt.Type {
	case "private":
		res, err = evt.API.Send(evt.Sender.UserID, msg, auto, async, amy.Private)
		break
	case "group":
		var msg message.CQGroup
		evt.ReadJSON(&msg)
		res, err = evt.API.Send(msg.GroupID, msg, auto, async, amy.Group)
		break
	case "discuss":
		var msg message.CQDiscuss
		evt.ReadJSON(&msg)
		res, err = evt.API.Send(msg.DiscussID, msg, auto, async, amy.Discuss)
		break
	}
	return
}
