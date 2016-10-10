package server

import (
	"log"
	"net/http"
	"golang.org/x/net/websocket"
)

type Server struct {
	pattern     string
	connections []*websocket.Conn
	clients     map[string]*Client
}

func NewServer(pattern string) *Server {
	connections := make([]*websocket.Conn, 10)
	clients := make(map[string]*Client)

	return &Server{
		pattern,
		connections,
		clients,
	}
}

func (s *Server) addClient(ws *websocket.Conn, msg *Message) {
	client := NewClient(ws, msg.Attribute("type"), msg.Attribute("id"))
	s.clients[client.Id] = client
}

func (s *Server) addConnection(connection *websocket.Conn) {
	s.connections = append(s.connections, connection)
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

		s.addConnection(ws)

		var msg Message
		err := websocket.JSON.Receive(ws, &msg)

		if err != nil {
			log.Fatal(err)
		} else {
			if msg.Attribute("type") == ACTION_AUTH {
				s.addClient(ws, msg)
			}
		}

	}
	http.Handle(s.pattern, websocket.Handler(onConnected))

	for true {}
}