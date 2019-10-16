package server

import (
	"encoding/json"
	"context"
	"net/http"
	"bytes"
)

// CQEvent 酷Q上报事件结构
type CQEvent struct {
	Type string
	writer http.ResponseWriter
	req *http.Request
	ctx context.Context
	handler []CQEventHandler
	hlength int
}

// CQEventHandler 上报数据处理函数
type CQEventHandler func(receive *CQEvent)

// newEvent 生成一个新的事件
func newEvent(ctx context.Context) CQEvent{
	var res CQEvent
	res.ctx = ctx
	wdata := res.ctx.Value(WriterKey)
	if w, ok := wdata.(http.ResponseWriter); ok{
		res.writer = w
	}
	rdata := res.ctx.Value(RequestKey)
	if r, ok := rdata.(*http.Request); ok {
		res.req = r
	}
	return res
}

func (event *CQEvent) write(data []byte) (int, error){	
	return event.writer.Write(data)
}
func (event *CQEvent) header() http.Header{
	return event.writer.Header()
}
func (event *CQEvent) writerHeader(statusCode int) {
	event.writer.WriteHeader(statusCode)
}
// JSON 快速响应JSON数据
func (event *CQEvent) JSON(data map[string]interface{}) (int, error){
	event.header().Set("Content-type", "application/json")
	event.writerHeader(200)
	bytes, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return event.write(bytes)
}

// context相关函数

// Value 获取值
func (event *CQEvent) Value(key interface{}) interface{}{
	return event.ctx.Value(key)
}
// Set 设置ctx值
func (event *CQEvent) Set(key interface{}, value interface{}) {
	event.ctx = context.WithValue(event.ctx, key, value)
}
// Next 调用下一个handler
func (event *CQEvent) Next() {
	if event.hlength == 0{
		return
	}
	handler := event.handler[0]
	event.handler = event.handler[1:]
	event.hlength = event.hlength - 1
	handler(event)
}

// Body 获取上报的原始信息
func (event *CQEvent) Body() []byte{
	return event.Value(ByteData).([]byte)
}
func (event *CQEvent) reqHeader() http.Header{
	return event.req.Header
}
// FormValue 获取表单字段
func (event *CQEvent) FormValue(key string) string{
	return event.req.FormValue(key)
}
// ReadJSON 将数据解析到msg中
func (event *CQEvent) ReadJSON(msg interface{}) error{
	data := event.Value(ByteData).([]byte)
	return json.Unmarshal(data, msg)
}
// MessageType 获取消息类型
// 只有在post_type类型为message时有用
func (event *CQEvent) MessageType() string {
	msg := event.Map()
	return msg["message_type"].(string)
}
// Map 获取map消息
func (event *CQEvent) Map() map[string]interface{} {	
	return loadintomap(event.Value(ByteData).([]byte))
}

func loadintomap(data []byte) (res map[string]interface{}) {
	decode := json.NewDecoder(bytes.NewReader(data))
	decode.UseNumber()
	decode.Decode(&res)
	return
}