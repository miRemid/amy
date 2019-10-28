package websocket

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"log"
)

type baseClient struct{
	url string
	port int
	conn *websocket.Conn
	token string
}

func (c *baseClient) Connect(t string){
	url := fmt.Sprintf("ws://%s:%d/%s/", c.url, c.port, t)
	var header http.Header
	if c.token != "" {
		header["Authorization"]	= []string{fmt.Sprintf("Bearer %s", c.token)}		
	}
	log.Println(url)
	conn,_,err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
}
