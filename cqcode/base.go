package cqcode

import (
	"fmt"
	"strings"
)

// CQParams cq码参数
type CQParams map[string]string

// CqCode 生成cqcode
func CqCode(function string, params CQParams) string {
	msg := fmt.Sprintf("[CQ:%s,", function)
	for k, v := range params{
		msg = msg + fmt.Sprintf("%s=%s,", k, v)
	}
	return strings.TrimRight(msg, ",") + "]"
}