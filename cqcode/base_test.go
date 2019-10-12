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