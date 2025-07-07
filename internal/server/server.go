package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/lealre/httpfromtcp/internal/request"
	"github.com/lealre/httpfromtcp/internal/response"
)

// Server is an HTTP 1.1 server
type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
		handler:  handler,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	resp := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		resp.WriteStatusLine(response.InternalServerError)
		body := []byte("error reading the request")
		resp.WriteHeaders(response.GetDefaultHeaders(len(body)))
		resp.WriteBody(body)
		return
	}
	s.handler(resp, req)
}
