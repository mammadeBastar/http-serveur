package server

import (
	"fmt"
	"http-serveur/internal/request"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	port     uint16
	closed   atomic.Bool
	listener net.Listener
}

func (s *Server) handle(conn net.Conn) {
	r, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatal("error", "error", err)
	}
	fmt.Printf("request line: \n")
	fmt.Printf("- Method: %s\n", r.RequestLine.Method)
	fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	fmt.Printf("headers: \n")
	r.Headers.ForEach(func(n, v string) {
		fmt.Printf("- %s: %s\n", n, v)
	})
	fmt.Printf("Body: \n")
	fmt.Println(string(r.Body))

	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain Content-Length: 13\r\n\r\nHello World!")
	conn.Write(out)
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

func Serve(port uint16) (*Server, error) {
	server := &Server{
		port:   port,
		closed: atomic.Bool{},
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
