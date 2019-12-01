package event

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/miRemid/amy"
)

type CQEventHandler func(evt CQEvent)

type CQEvent interface {
	JSON(statuscode int, data interface{}) (int, error)
	ReqHeader() http.Header
	Next()
	GetBody() []byte
}

type Map map[string]interface{}

type CQEventBase struct {
	writer   http.ResponseWriter
	Request  *http.Request
	httpFlag bool
	handler  []CQEventHandler
	flag     bool

	API  *amy.API
	Body []byte
}

func (evt *CQEventBase) Use(handlers ...CQEventHandler) {
	for _, handler := range handlers {
		evt.handler = append(evt.handler, handler)
	}
}

func (evt *CQEventBase) SetHTTP(w http.ResponseWriter, r *http.Request) {
	if !evt.httpFlag {
		evt.writer = w
		evt.Request = r
		evt.httpFlag = true
	}
}

func (evt CQEventBase) ReqHeader() http.Header {
	return evt.Request.Header
}

func (evt CQEventBase) Next() {
	if len(evt.handler) == 0 {
		return
	}
	handler := evt.handler[0]
	evt.handler = evt.handler[1:]
	handler(evt)
}

func (evt CQEventBase) GetBody() []byte {
	return evt.Body
}

func (evt CQEventBase) JSON(statuscode int, data interface{}) (int, error) {
	if evt.flag {
		return 0, fmt.Errorf("amy http response error: already response")
	}
	evt.flag = true
	evt.writer.Header().Set("Content-type", "application/json")
	evt.writer.WriteHeader(statuscode)
	if data == nil {
		return 0, nil
	}
	bytedata, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return evt.writer.Write(bytedata)
}
