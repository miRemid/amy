package utils

import "strings"

func CmdParser(message string, cmds ...string) (cmd, params string) {
	msg := strings.TrimSpace(message)
	split := strings.Split(msg, " ")
	tcmd := split[0]
	if tcmd == "" {
		return "", ""
	}
	for _, v := range cmds {
		if len(v) == len(tcmd) && v != tcmd || len(v) > len(tcmd) {
			continue
		}
		if v == tcmd[:len(v)] {
			if len(tcmd[len(v):]) == 0 {
				continue
			}
			cmd = tcmd[len(v):]
			if len(split) == 1 {
				params = ""
			} else {
				params = msg[len(cmd)+1:]
			}
			break
		}
	}
	return
}