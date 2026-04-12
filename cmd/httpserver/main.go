package main

import (
	"http-serveur/internal/headers"
	"http-serveur/internal/reponse"
	"http-serveur/internal/request"
	"http-serveur/internal/server"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const port = 42069

const (
	YourProblemBody = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
	MyProblemBody = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
	OkBody = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
)

func main() {
	var handle server.Handler = func(w reponse.Writer, req *request.Request) {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			err := w.WriteStatusLine(400)
			if err != nil {
				log.Fatal("error writing status line", "error", err)
			}
			h := headers.NewHeaders()
			var b []byte = []byte(YourProblemBody)
			h.Set("Content-Length", strconv.Itoa(len(b)))
			h.Set("Connection", "close")
			h.Set("Content-Type", "text/html")
			err = w.WriteHeaders(*h)
			if err != nil {
				log.Fatal("error writing headers", "error", err)
			}

			_, err = w.WriteBody(b)
			if err != nil {
				log.Fatal("error writing status line", "error", err)
			}

		} else if req.RequestLine.RequestTarget == "/myproblem" {
			err := w.WriteStatusLine(500)
			if err != nil {
				log.Fatal("error writing status line", "error", err)
			}
			h := headers.NewHeaders()
			var b []byte = []byte(MyProblemBody)
			h.Set("Content-Length", strconv.Itoa(len(b)))
			h.Set("Connection", "close")
			h.Set("Content-Type", "text/html")
			err = w.WriteHeaders(*h)
			if err != nil {
				log.Fatal("error writing headers", "error", err)
			}
			_, err = w.WriteBody(b)
			if err != nil {
				log.Fatal("error writing status line", "error", err)
			}
		} else {
			err := w.WriteStatusLine(200)
			if err != nil {
				log.Fatal("error writing status line", "error", err)
			}
			h := headers.NewHeaders()
			var b []byte = []byte(OkBody)
			h.Set("Content-Length", strconv.Itoa(len(b)))
			h.Set("Connection", "close")
			h.Set("Content-Type", "text/html")
			err = w.WriteHeaders(*h)
			if err != nil {
				log.Fatal("error writing headers", "error", err)
			}
			_, err = w.WriteBody(b)
			if err != nil {
				log.Fatal("error writing status line", "error", err)
			}
		}

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
