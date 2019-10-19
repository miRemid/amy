package amy

const (
	// Private 私人消息
	Private = iota
	// Group 群组消息
	Group
	// Discuss 讨论组消息
	Discuss
)

// Response 返回消息结构体
type Response struct {
	Reader []byte
	Error	error
}
// CResponse 返回管道
type CResponse chan Response