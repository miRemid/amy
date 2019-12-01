package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CQEvent 酷Q上报事件结构
type CQEvent struct {
	Type    string
	writer  http.ResponseWriter
	req     *http.Request
	ctx     context.Context
	handler []CQEventHandler
	hlength int
	flag    bool
}

// String 打印消息
func (event CQEvent) String() string {
	return string(event.Value(ByteData).([]byte))
}

// CQEventHandler 上报数据处理函数
type CQEventHandler func(receive CQEvent)

func (event CQEvent) write(data []byte) (int, error) {
	return event.writer.Write(data)
}
func (event CQEvent) header() http.Header {
	return event.writer.Header()
}
func (event CQEvent) writerHeader(statusCode int) {
	event.writer.WriteHeader(statusCode)
}

// JSON 快速响应JSON数据
func (event CQEvent) JSON(statuscode int, data map[string]interface{}) (int, error) {
	if event.flag {
		return 0, fmt.Errorf("amy http response error: already response")
	}
	event.flag = true
	event.header().Set("Content-type", "application/json")
	event.writerHeader(statuscode)
	bytes, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return event.write(bytes)
}

// Value 获取值
func (event CQEvent) Value(key interface{}) interface{} {
	return event.ctx.Value(key)
}

// Set 设置ctx值
func (event *CQEvent) Set(key interface{}, value interface{}) {
	event.ctx = context.WithValue(event.ctx, key, value)
}

// Next 调用下一个handler
func (event CQEvent) Next() {
	if event.hlength == 0 {
		return
	}
	handler := event.handler[0]
	event.handler = event.handler[1:]
	event.hlength = event.hlength - 1
	handler(event)
}

// Body 获取上报的原始信息
func (event *CQEvent) Body() []byte {
	return event.Value(ByteData).([]byte)
}
func (event *CQEvent) reqHeader() http.Header {
	return event.req.Header
}

// ReadJSON 将数据解析到msg中
func (event CQEvent) ReadJSON(msg interface{}) error {
	data := event.Value(ByteData).([]byte)
	return json.Unmarshal(data, msg)
}

// MessageType 获取消息类型
// 只有在post_type类型为message时有用
func (event CQEvent) MessageType() string {
	msg := event.Map()
	return msg["message_type"].(string)
}

// Map 获取map消息
func (event CQEvent) Map() map[string]interface{} {
	return loadintomap(event.Value(ByteData).([]byte))
}

func loadintomap(data []byte) (res map[string]interface{}) {
	decode := json.NewDecoder(bytes.NewReader(data))
	decode.UseNumber()
	decode.Decode(&res)
	return
}
