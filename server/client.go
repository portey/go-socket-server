package server

import (
	"golang.org/x/net/websocket"
	"io"
	"fmt"
)

const channelBufSize = 100

const CLIENT_TYPE_ACCESS_POINT string = "LOCK"
const CLIENT_TYPE_USER string = "USER"
const CLIENT_TYPE_CLOUD string = "CLOUD"

type Client struct {
	Id         string
	ClientType string
	server     *Server
	ws         *websocket.Conn
	messages   chan *Message
	exit       chan bool
}

func NewClient(ws *websocket.Conn, server *Server, clientType string, id string) *Client {
	messages := make(chan *Message, channelBufSize)
	exit := make(chan bool)

	return &Client{id, clientType, server, ws, messages, exit, }
}

func (c *Client) Write(msg *Message) {
	select {
	case c.messages <- msg:
	default:
		err := fmt.Errorf("client %d is disconnected.", c.Id)
		c.server.Err(err)

		c.server.removeClient(c)
	}
}

func (c *Client) Listen() {
	go c.listenReads()
	c.listenWrites()
}

func (c *Client) listenReads() {
	for {
		select {

		case <-c.exit:
			c.exit <- true
			return

		default:
			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.exit <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.server.onMessage(c.ws, &msg)
			}
		}
	}

}

func (c *Client) listenWrites() {
	for {
		select {

		case msg := <-c.messages:
			websocket.JSON.Send(c.ws, msg)

		case <-c.exit:
			c.exit <- true
			return
		}
	}
}

