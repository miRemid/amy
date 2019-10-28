package model

import (
	"github.com/miRemid/amy/utils"
)

// CQParams 数据参数
type CQParams map[string]interface{}

// CQMessage WebSocket发送的消息
type CQMessage struct{
	API string			`json:"action"`
	Params CQParams		`json:"params"`
}

// CQResponse 回复消息
type CQResponse struct {
	Status	string		`json:"status"`
	RetCode	int			`json:"ret_code"`
	Data	interface{}	`json:"data"`
}

// NewCQResponse return a CQResponse
func NewCQResponse(data []byte) (CQResponse, error){
	var res CQResponse
	err := utils.LoadIntoStruct(data, &res)
	return res, err
}