package message

import (
	"fmt"
)

// build.go 包含了消息创建器结构
// 所有的消息需要通过一个Builder来创建

// CQMsgBuilder 消息序列化器
type CQMsgBuilder struct {
}

// NewCQMsgBuilder 创建新的序列化器
func NewCQMsgBuilder() *CQMsgBuilder {
	return &CQMsgBuilder{}
}

// PrivateMsg 创建一条私人消息
// to 目标id
// message 消息
// auto 消息是否为字符串
func (builder *CQMsgBuilder) PrivateMsg(to int, message interface{}, auto bool) CQPrivateMsg{
	return CQPrivateMsg{
		UserID: to,
		CQSendMsgBase: CQSendMsgBase{
			Message: message,
			Auto:	auto,
		},
	}
}

// GroupMsg 创建一条群消息
// to 目标id
// message 消息
// auto 消息是否为字符串
func (builder *CQMsgBuilder) GroupMsg(to int, message interface{}, auto bool) CQGroupMsg{
	return CQGroupMsg{
		GroupID: to,
		CQSendMsgBase: CQSendMsgBase{
			Message: message,
			Auto:	auto,
		},
	}
}

// DiscussMsg 创建一条讨论组消息
// to 目标id
// message 消息
// auto 消息是否为字符串
func (builder *CQMsgBuilder) DiscussMsg(to int, message interface{}, auto bool) CQDiscussMsg{
	return CQDiscussMsg{
		DiscussID: to,
		CQSendMsgBase: CQSendMsgBase{
			Message: message,
			Auto:	auto,
		},
	}
}

// RawMsg 创建一条原生消息
// msgType 消息类型
// to 目标id
// message 消息
// auto 消息是否为字符串
func (builder *CQMsgBuilder) RawMsg(msgType string, to int, message interface{}, auto bool) CQRawMsg{
	return CQRawMsg{
		MessageType: msgType,
		UserID:	to,
		GroupID: to,
		DiscussID: to,		
		CQSendMsgBase: CQSendMsgBase{
			Message: message,
			Auto:	auto,
		},
	}
}

// CQJSON 生成一条JSON格式消息段
// t 消息块信息类型
// kv 消息块数据键值对
func (builder *CQMsgBuilder) CQJSON(t string, kv ...string) CQJSON {
	var res CQJSON
	if len(kv) % 2 != 0 {
		return res
	}
	res.Type = t
	data := make(map[string]interface{})
	for {
		if len(kv) == 0{
			break
		}
		slice := kv[0:2]
		fmt.Println(slice)
		data[slice[0]] = slice[1]
		kv = kv[2:]
	}
	res.Data = data
	return res
}