package server

import (
	"golang.org/x/net/websocket"
)

const channelBufSize = 100

const CLIENT_TYPE_ACCESS_POINT string = "LOCK"
const CLIENT_TYPE_USER string = "USER"
const CLIENT_TYPE_CLOUD string = "CLOUD"

type Client struct {
	Id         string
	ClientType string
	connection *websocket.Conn
}

func NewClient(connection *websocket.Conn, clientType string, id string) *Client {
	return &Client{id, clientType, connection, }
}