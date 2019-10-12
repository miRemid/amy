package amy

// Response 返回消息结构体
type Response struct {
	Reader []byte
	Error	error
}
// CResponse 返回管道
type CResponse chan Response