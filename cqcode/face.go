package cqcode

import (
	"fmt"
)

// face.go 包含了qq所有自带表情包数字字段

const (

)

// Face 返回自带表情cqcode
func Face(id int) string {
	return fmt.Sprintf("[CQ:face,id=%d]", id)
}

// BFace 返回原创表情cqcode
// 原创表情存放在酷Q目录的data\bface\下
func BFace(id int) string {
	return fmt.Sprintf("[CQ:bface,id=%d]", id)
}

// SFace 小表情
func SFace(id int) string {
	return fmt.Sprintf("[CQ:sface,id=%d]", id)
}