package server

import (
	"fmt"
	"http-serveur/internal/reponse"
	"http-serveur/internal/request"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	port     uint16
	closed   atomic.Bool
	listener net.Listener
	handler  Handler
}

type HandlerError struct {
	StatusCode reponse.StatusCode
	Msg        string
}

type Handler func(w reponse.Writer, req *request.Request)

func (s *Server) handle(conn net.Conn) {
	r, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatal("error", "error", err)
	}

	writer := reponse.NewWriter(conn)

	s.handler(*writer, r)

	conn.Close()
}

func (s *Server) listen() {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if s.closed.Load() {
				return
			}
			if err != nil {
				log.Fatal("error", "error", err)
			}
			go s.handle(conn)
		}
	}()
}

func writeError(w *reponse.Writer, err *HandlerError) error {
	h := reponse.GetDefaultHeaders(len(err.Msg), "text/plain")
	w.WriteStatusLine(err.StatusCode)
	w.WriteHeaders(h)
	w.Write([]byte(err.Msg))
	return nil
}

func Serve(port uint16, handler Handler) (*Server, error) {
	server := &Server{
		port:    port,
		closed:  atomic.Bool{},
		handler: handler,
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.port))
	server.listener = listener
	if err != nil {
		return nil, err
	}

	server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	s.listener.Close()
	return nil
}
