package websocket

import (
	"net/http"
	"github.com/gorilla/websocket"	
	"log"
)

type baseClient struct{
	url string
	port int
	conn *websocket.Conn
	token string
}

func (c *baseClient) Connect(url string, header http.Header){	
	conn,_,err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
}
