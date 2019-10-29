package websocket

import (
	"log"
	"fmt"
	"sync"
	"net/http"
	"github.com/miRemid/amy/websocket/model"
)

// APIHandler a
type APIHandler func(res model.CQResponse)

// APIClient recall the CQHTTP API websocket
type APIClient struct {
	baseClient
	handler APIHandler
}

// NewAPIClient return a new APIClient
func NewAPIClient(url string, port int, token string) *APIClient{
	var client APIClient
	client.URL = url
	client.Port = port
	client.token = token
	client.handler = func(res model.CQResponse){
		log.Printf(res.Status)
	}
	return &client
}

// Connect to the ws server
func (c *APIClient) Connect() {
	url := fmt.Sprintf("ws://%s:%d/api/", c.URL, c.Port)
	var header http.Header
	if c.token != "" {
		header["Authorization"]	= []string{fmt.Sprintf("Bearer %s", c.token)}		
	}
	c.baseClient.Connect(url, header)
}

// SetToken set your access_token
func (c *APIClient) SetToken(token string) {
	c.token = token
}

func (c *APIClient) receive(wg *sync.WaitGroup) {
	for {
		_, body, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("read error:%v", err)
			continue
		}
		if res, err := model.NewCQResponse(body); err != nil {
			log.Printf("parse error:%v", err)
		}else{
			c.handler(res)
		}
		break
	}
	wg.Done()
}

// Send 发送消息
func (c *APIClient) Send(api string, params model.CQParams) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	c.Connect()
	go c.receive(&wg)
	go func(wg *sync.WaitGroup){
		msg := model.CQMessage{
			API: api,
			Params: params,
		}
		c.conn.WriteJSON(&msg)
		wg.Done()
	}(&wg)	
	wg.Wait()
}

// OnResponse set the handler function
func (c *APIClient) OnResponse(handler APIHandler) {
	c.handler = handler
}