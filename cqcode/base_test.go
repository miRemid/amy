package cqcode

import (
	"testing"
)

func TestCqCode(t *testing.T) {
	msg := CqCode("test", CQParams{
		"file": "asdf",
	})
	t.Log(msg)
}

func TestCQParse(t *testing.T) {
	cqstr := "[CQ:at,id=12323534]"
	cqcode := CQParse(cqstr)
	t.Log(cqcode)
}
