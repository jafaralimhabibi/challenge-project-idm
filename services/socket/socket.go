package services

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

var serv *socketio.Server

// NewServer creates and starts a socket.io server (minimal handlers).
func NewServer() *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("socket connected:", s.ID())
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("socket disconnected:", s.ID(), "reason:", reason)
	})

	// run server goroutine
	go server.Serve()

	// set global
	SetServer(server)
	return server
}

func SetServer(s *socketio.Server) { serv = s }
func GetServer() *socketio.Server  { return serv }

// EmitToNamespace broadcasts to namespace; safe no-op if server nil.
func EmitToNamespace(namespace, event string, v interface{}) {
	if serv == nil {
		return
	}
	serv.BroadcastToNamespace(namespace, event, v)
}
