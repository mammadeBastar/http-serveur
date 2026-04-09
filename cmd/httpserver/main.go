package main

import (
	"http-serveur/internal/request"
	"http-serveur/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	var handle server.Handler = func(w io.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			return &server.HandlerError{
				StatusCode: 400,
				Msg:        "Your problem is not my problem\n",
			}
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			return &server.HandlerError{
				StatusCode: 500,
				Msg:        "Woopsie, my bad\n",
			}
		}

		w.Write([]byte("All good, frfr\n"))

		return nil
	}
	s, err := server.Serve(uint16(port), handle)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

