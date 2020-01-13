package utils

import "testing"

import "log"

func TestCmdParser(t *testing.T) {
	cmd, params := CmdParser("!setid", "!", "！")
	log.Println(cmd, params)
	for _, v := range params {
		log.Println(v)
	}
}
