package websocket

import (
	"github.com/miRemid/amy/websocket/model"
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
	"log"
)

type baseClient struct{
	URL string
	Port int
	conn *websocket.Conn
	token string
}

// Connect to the cqhttp
func (c *baseClient) Connect(url string, header http.Header){	
	conn,_,err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
}


// CQClient is a websocket client connect to the cqhttp
type CQClient struct{
	MsgClient
	pool *sync.Pool
	apiurl string
	apiport int
	apihandler APIHandler
}

// NewCQClient retuan a new cqclient
// if no token input ""
func NewCQClient(url string, port int) *CQClient{
	var res CQClient
	res.URL = url
	res.apiurl = url
	res.Port = port
	res.apiport = port
	return &res
}

// OnResponse se the handler function of response
func (c *CQClient) OnResponse(handler APIHandler) {
	c.apihandler = handler
}

// Send a message
func (c *CQClient) Send(api string, params model.CQParams) {
	client := c.pool.Get().(*APIClient)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go client.receive(&wg)
	go func(c *APIClient, wg *sync.WaitGroup) {
		msg := model.CQMessage{
			API: api,
			Params: params,
		}
		client.conn.WriteJSON(&msg)
	}(client, &wg)
	wg.Wait()
	c.pool.Put(client)
}

// Run the client
func (c *CQClient) Run() {
	c.pool = &sync.Pool{
		New: func() interface{}{
			var api APIClient
			api.URL = c.apiurl
			api.Port = c.apiport
			if c.token != "" {
				api.token = c.token
			}
			api.Connect()
			api.handler = c.apihandler		
			return &api
		},
	}
	c.MsgClient.Run()
}

// SetAPIConfig will set apiclient config
func (c *CQClient) SetAPIConfig(url string, port int) {
	c.apiurl = url
	c.apiport = port
}

// SetToken will set the apiclient's token
func (c *CQClient) SetToken(token string) {
	c.token = token
}