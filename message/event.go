package message

import (
	"time"
)

// cqmsg.go 包含接受的消息格式

// CQEventStruct 空接口
type CQEventStruct interface{}

// CQReceive 基类，每条消息必有的消息结构
type CQReceive struct {
	PostType	string			`json:"post_type"`
	Time		time.Duration	`json:"time"`
	SelfID		int				`json:"self_id"`
}

// CQPrivate 私聊消息
type CQPrivate struct {
	CQReceive
	MessageType	string			`json:"message_type"`
	SubType		string			`json:"sub_type"`
	MessageID	int32			`json:"message_id"`
	UserID		int				`json:"user_id"`
	Message		string			`json:"message"`
	RawMessage	string			`json:"raw_message"`
	Font		int32			`json:"font"`
	Sender		CQSender		`json:"sender"`	
}

// CQGroup 群聊消息
type CQGroup struct {
	CQReceive
	MessageType	string			`json:"message_type"`
	SubType		string			`json:"sub_type"`
	MessageID	int32			`json:"message_id"`
	GroupID		int64			`json:"group_id"`
	UserID		int				`json:"user_id"`
	Anonymous	CQAnonymous		`json:"anonymous"`
	Message		string			`json:"message"`	
	RawMessage	string			`json:"raw_message"`
	Font		int32			`json:"font"`
	Sender		CQGroupSender	`json:"sender"`
}

// CQDiscuss 讨论组消息
type CQDiscuss struct {
	CQReceive
	MessageType	string			`json:"message_type"`
	MessageID	int32			`json:"message_id"`
	DiscussID	int				`json:"discuss_id"`
	UserID		int				`json:"user_id"`
	Message		string 			`json:"message"`	
	RawMessage	string			`json:"raw_message"`
	Font		int32			`json:"font"`
	Sender		CQSender		`json:"sender"`	
}

// CQGroupFile 群文件上传消息
type CQGroupFile struct {
	CQReceive
	NoticeType	string	`json:"notice_type"`
	GroupID		int		`json:"group_id"`
	UserID		int		`json:"user_id"`
	File		CQFile	`json:"file"`
}

// CQAdminChange 群管理员变动
type CQAdminChange struct {
	CQReceive
	SubType		string	`json:"sub_type"`
	NoticeType	string	`json:"notice_type"`
	GroupID		int		`json:"group_id"`
	UserID		int		`json:"user_id"`
}

// CQGuDecrease 群成员减少
type CQGuDecrease struct {
	CQReceive
	SubType		string	`json:"sub_type"`
	NoticeType	string	`json:"notice_type"`
	GroupID		int		`json:"group_id"`
	UserID		int		`json:"user_id"`
	OperatorID	int		`json:"operator_id"`
}

// CQGuAdd 群成员增加
// 相应消息和群成员减少一致
type CQGuAdd struct {
	CQGuDecrease
}

// CQFriendAdd 好友添加
type CQFriendAdd struct {
	CQReceive
	NoticeType	string	`json:"notice_type"`
	UserID		int		`json:"user_id"`
}

// CQFriendRequest 好友请求
type CQFriendRequest struct {
	CQReceive
	RequestType string	`json:"request_type"`
	UserID		int		`json:"user_id"`
	Comment		string	`json:"comment"`
	Flag		string	`json:"flag"`
}

// CQGroupRequest 加群请求/邀请
type CQGroupRequest struct {
	CQReceive
	RequestType string	`json:"request_type"`
	SubType		string	`json:"sub_type"`
	GroupID		int		`json:"group_id"`
	UserID		int		`json:"user_id"`
	Comment		string	`json:"comment"`
	Flag 		string	`json:"flag"`
}