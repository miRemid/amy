package cqcode

import (
	"fmt"
)

// At @某人
func At(userid int) string {
	return fmt.Sprintf("[CQ:at,qq=%d]", userid)
}

// AtAll @全体成员
func AtAll() string {
	return fmt.Sprintf("[CQ:at,qq=all]")
}

// RPS 猜拳魔法表情
func RPS() string {
	return fmt.Sprintf("[CQ:rps,type=1]")
}

// Dice 掷骰子表情
func Dice() string {
	return fmt.Sprintf("[CQ:dice,type=1]")
}

// Shake 戳一戳，仅好友消息
func Shake() string {
	return fmt.Sprintf("[CQ:shake]")	
}