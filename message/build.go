package message

// build.go 包含了消息创建器结构
// 所有的消息需要通过一个Builder来创建

// CQMsgBuilder 消息序列化器
type CQMsgBuilder struct {
	// 消息格式类型，可选string或者array
	MsgType	string
}

// NewCQMsgBuilder 创建新的序列化器
// 消息类型默认为字符串
func NewCQMsgBuilder() *CQMsgBuilder {
	return &CQMsgBuilder{
		MsgType: "string",
	}
}

// PrivateMsg 创建一条私人消息
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