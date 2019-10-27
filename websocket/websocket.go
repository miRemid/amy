package websocket

import (
	"log"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/miRemid/amy/websocket/model"
)

// Client is the struct of cq websocket client
type Client struct{
	url string
	port int
	conn *websocket.Conn
	message WebEventHandler
}

// WebEventHandler is the handler function of receive message
type WebEventHandler func(event model.CQEvent)

// NewClient return a new websocket server ptr
func NewClient(url string, port int) *Client{
	return &Client{
		url : url,
		port: port,
	}
}

func (c *Client) connect(){
	url := fmt.Sprintf("ws://%s:%d/event/", c.url, c.port)
	log.Println(url)
	conn,_,err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
}

// Run will open a websocket client
func (c *Client) Run() {
	c.connect()	
	for {
		_, body, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("read error:%v", err)
		}
		evt := model.NewCQEvent(body)
		go c.message(evt)
	}
}

// OnMessage will set the message parse function
func (c *Client) OnMessage(handler WebEventHandler){
	c.message = handler
}

// Demo will show a demo websocket client receive the data from CQHTTP webscoket server
func Demo() {
	client := NewClient("127.0.0.1", 6700)
	client.OnMessage(func(evt model.CQEvent){
		log.Println(evt.Type)
	})
	client.Run()
}