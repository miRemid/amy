package cqcode

import (
	"fmt"
	"strings"
)

// CQCode CQ码
type CQCode struct {
	Func   string
	Params CQParams
}

// CQParams cq码参数
type CQParams map[string]string

func (c CQCode) String() string {
	msg := fmt.Sprintf("[CQ:%s,", c.Func)
	for k, v := range c.Params {
		msg = msg + fmt.Sprintf("%s=%s,", k, v)
	}
	return strings.TrimRight(msg, ",") + "]"
}

// CqCode 生成cqcode
func CqCode(function string, params CQParams) string {
	msg := fmt.Sprintf("[CQ:%s,", function)
	for k, v := range params {
		msg = msg + fmt.Sprintf("%s=%s,", k, v)
	}
	return strings.TrimRight(msg, ",") + "]"
}

// CQParse 解析CQ码字符串
func CQParse(cqstr string) (code CQCode) {
	code.Params = make(map[string]string)
	str := cqstr[1 : len(cqstr)-1]
	list := strings.Split(str, ",")
	code.Func = strings.Split(list[0], ":")[1]
	for _, v := range list[1:] {
		params := strings.Split(v, "=")
		code.Params[string(params[0])] = string(params[1])
	}
	return
}
