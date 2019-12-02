package event

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/miRemid/amy"
)

// CQEventHandler is the base event middleware function
type CQEventHandler func(evt CQEvent)

// Map is the event json data struct
type Map map[string]interface{}

// CQEvent is the base CQHTTP event struct
type CQEvent struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	httpFlag bool
	handler  []CQEventHandler
	flag     bool

	API  *amy.API
	Body []byte
}

// Use the handlers
func (evt *CQEvent) Use(handlers ...CQEventHandler) {
	for _, handler := range handlers {
		evt.handler = append(evt.handler, handler)
	}
}

// Next to the next handler function
func (evt CQEvent) Next() {
	if len(evt.handler) == 0 {
		return
	}
	handler := evt.handler[0]
	evt.handler = evt.handler[1:]
	handler(evt)
}

// JSON reply the CQHTTP server
func (evt CQEvent) JSON(statuscode int, data interface{}) (int, error) {
	if evt.flag {
		return 0, fmt.Errorf("amy http response error: already response")
	}
	evt.flag = true
	evt.Writer.Header().Set("Content-type", "application/json")
	evt.Writer.WriteHeader(statuscode)
	if data == nil {
		return 0, nil
	}
	bytedata, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return evt.Writer.Write(bytedata)
}
