package message

// send.go 包含了发送消息的结构体

import (
	"encoding/json"
	"bytes"
)

// CQMessage 发送消息接口
type CQMessage interface{
	// 返回http可发送的数据
	Value() (*bytes.Reader, error)
}

// CQMAP map
type CQMAP map[string]interface{}

// CQJSON 消息段
type CQJSON struct{
	Type string					`json:"type"`
	Data map[string]interface{}	`json:"data"`
}

// CQArray array格式数据
type CQArray []CQJSON

// Value 返回http可发送的消息
func (m CQMAP) Value() (*bytes.Reader, error) {
	byteData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(byteData), nil
}

// CQSendMsgBase 每条发送的消息通用的字段
type CQSendMsgBase struct {
	Message interface{} `json:"message"`
	Auto	bool		`json:"auto_escape"`
}

// CQPrivateMsg 发送私聊消息
type CQPrivateMsg struct {
	CQSendMsgBase
	UserID	int			`json:"user_id"`	
}

// Value 返回http可发送的消息
func (m CQPrivateMsg) Value() (*bytes.Reader, error) {
	byteData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(byteData), nil
}

// CQGroupMsg 发送群聊消息
type CQGroupMsg struct {
	CQSendMsgBase
	GroupID	int			`json:"group_id"`	
}

// Value 返回http可发送的消息
func (m CQGroupMsg) Value() (*bytes.Reader, error) {
	byteData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(byteData), nil
}

// CQDiscussMsg 发送讨论组消息
type CQDiscussMsg struct {
	CQSendMsgBase
	DiscussID	int		`json:"discuss_id"`	
}

// Value 返回http可发送的消息
func (m CQDiscussMsg) Value() (*bytes.Reader, error) {
	byteData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(byteData), nil
}

// CQRawMsg 发送自定义消息
type CQRawMsg struct {
	CQSendMsgBase
	MessageType string	`json:"message_type"`
	DiscussID	int		`json:"discuss_id"`	
	GroupID	int			`json:"group_id"`	
	UserID	int			`json:"user_id"`
}

// Value 返回http可发送的消息
func (m CQRawMsg) Value() (*bytes.Reader, error) {
	byteData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(byteData), nil
}