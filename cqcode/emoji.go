package cqcode

import (
	"fmt"
)

// Emoji 返回emoji的cqcode
func Emoji(unicode int) string {
	return fmt.Sprintf("[CQ:emoji,id=%d]", unicode)
}