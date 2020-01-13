package utils

import "strings"

// CmdParser 解析Cmd命令
func CmdParser(message string, cmds ...string) (cmd string, params []string) {
	msg := strings.TrimSpace(message)
	split := strings.Split(msg, " ")
	tcmd := split[0]
	if tcmd == "" {
		return "", nil
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
				params = nil
			} else {
				params = split[1:]
			}
			break
		}
	}
	return
}
