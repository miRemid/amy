package amy

import (
	"time"
	"net/http"
)

// API CQHTTP消息发送API
type API struct {
	CQHTTP	string
	CQPORT	int
	Client	*http.Client
	Token	string
}

// NewAmyAPI 创建Api实例链接CQ
func NewAmyAPI(cqhttp string, port int) *API{
	api := API{
		CQHTTP:	cqhttp,
		CQPORT: port,		
	}
	api.Client = &http.Client{
		Timeout: time.Second * 5,
	}
	return &api
}

// SetTimeout 设置api请求发送超时时间
func (api *API) SetTimeout(timeout time.Duration) {
	api.Client.Timeout = timeout
}

// SetToken 设置Access_Token
func (api *API) SetToken(token string) {
	api.Token = token
}