package message

const (
	// CQNeedToken 未提供access token
	CQNeedToken = 401
	// CQTokenNotCompared 提供的access token不符合
	CQTokenNotCompared = 403
	// CQContentTypeNotSupport Post请求的Content-Type不支持
	CQContentTypeNotSupport = 406
	// CQPostStructError Post请求正文格式不正确
	CQPostStructError = 400
	// CQAPINotFound API不存在
	CQAPINotFound = 404
	// CQStatusOK 操作成功或其他原因失败
	CQStatusOK = 200

	// RetCodeOK 操作成功
	RetCodeOK = 0
	// RetAsync 进入异步执行
	RetAsync = 1
	// RetParamInvalid 参数缺失或参数无效
	RetParamInvalid = 100
	// RetResponseInvalid 返回数据无效
	RetResponseInvalid = 102
	// RetOprationFailed 操作失败
	RetOprationFailed = 103
	// RetCookieOrCSRFError 酷Q 提供的凭证（Cookie 和 CSRF Token）失效
	RetCookieOrCSRFError = 104
	// RetInitError 工作线程池未正确初始化
	RetInitError = 201
)

var (
	// DefaultBuilder 默认构造器
	DefaultBuilder *CQMsgBuilder
)

func init() {
	DefaultBuilder = NewCQMsgBuilder()
}

// CQSender 消息发送方信息
type CQSender struct {
	// UserID 用户id
	UserID		int		`json:"user_id"`
	// NickName 用户昵称
	NickName	string	`json:"nickname"`
	// Sex 性别
	Sex			string	`json:"sex"`
	// Age 年龄
	Age			int32	`json:"age"`
}

// CQGroupSender 群消息发送方信息
type CQGroupSender struct {
	CQSender
	Card		string	`json:"card"`
	Area		string	`json:"area"`
	Level		string	`json:"level"`
	Role		string	`json:"role"`
	Title		string	`json:"title"`
}

// CQAnonymous 匿名信息
type CQAnonymous struct {
	ID			int64	`json:"id"`
	Name		string	`json:"string"`
	Flag		string	`json:"flag"`
}

// CQFile 文件信息
type CQFile struct {
	ID		string	`json:"id"`
	Name	string	`json:"name"`
	Size	int		`json:"size"`
	Busid	int		`json:"busid"`
}

// CQMsgSegment 消息段
type CQMsgSegment struct {
	MessageType	string					`json:"type"`
	Data		map[string]interface{}	`json:"data"`
}

// CQMessageID 发送的消息的ID
type CQMessageID struct {
	ID	int		`json:"message_id"`
}