package cqcode

import (
	"testing"
	"log"
)

func TestCqCode(t *testing.T) {
	msg := CqCode("test", CQParams{
		"file": "asdf",
	})
	t.Log(msg)
}

func TestCQParse(t *testing.T) {
	cqstr := "asdf[CQ:at,id=fdfsdf]asdf[CQ:at,id=asdfdf]dasfioghfdguihn[CQ:at,id=ghjgfhj]"
	cqcodes := CQSplit(cqstr)
	for _, v := range cqcodes {
		log.Printf("func=%v,params=%v\n", v.Func, v.Params)
	}
	t.Log(cqcodes)
}
