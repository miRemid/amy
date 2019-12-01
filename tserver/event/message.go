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

type CQSession struct {
	*CQEventBase

	Type       string
	Sender     message.CQSender
	Message    string
	RawMessage string
}

func (evt CQSession) Params(cmds ...string) (cmd, params string) {
	msg := string(evt.Body)
	return utils.CmdParser(msg, cmds...)
}

func (evt CQSession) CQCode() (res []cqcode.CQCode) {
	res = make([]cqcode.CQCode, 0)
	str := reg.FindAllString(evt.Message, -1)
	for _, v := range str {
		res = append(res, cqcode.CQParse(v))
	}
	return
}

func (evt CQSession) ReadJSON(cq interface{}) error {
	body := evt.Body
	return json.Unmarshal(body, cq)
}

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
