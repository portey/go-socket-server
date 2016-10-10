package server

import (
	"log"
	"net/http"
	"golang.org/x/net/websocket"
)

type Server struct {
	pattern string
	clients map[string]*Client
}

func NewServer(pattern string) *Server {
	clients := make(map[string]*Client)

	return &Server{
		pattern,
		clients,
	}
}

func (s *Server) Err(err error) {
	log.Println("Error:", err.Error())
}

func (s *Server) addNewClient(ws *websocket.Conn, msg *Message) {
	messageType, err := msg.Attribute("type");
	if err == nil {
		id, err := msg.Attribute("type");

		if err == nil {
			client := NewClient(ws, s, messageType, id)
			s.clients[client.Id] = client
			client.Listen()
		}
	}
}

func (s *Server) removeClient(client *Client) {
	delete(s.clients, client.Id)
}

func (s *Server) onMessage(ws *websocket.Conn, msg *Message) {
	messageType, err := msg.Attribute("type");

	if err == nil {
		if messageType == ACTION_AUTH {
			s.addNewClient(ws, msg)
		}
	}

}

func (s *Server) Listen() {

	log.Println("Listening server...")

	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		var msg Message
		err := websocket.JSON.Receive(ws, &msg)

		if err != nil {
			log.Fatal(err)
		} else {
			s.onMessage(ws, &msg)

		}

	}
	http.Handle(s.pattern, websocket.Handler(onConnected))

	for true {}
}