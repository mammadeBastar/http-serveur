package main

import (
	"errors"
	"fmt"
	"http-serveur/internal/headers"
	"http-serveur/internal/reponse"
	"http-serveur/internal/request"
	"http-serveur/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
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

var clinet = &http.Client{
	Timeout: 30 * time.Second,
}

func main() {
	var handle server.Handler = func(w reponse.Writer, req *request.Request) {
		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			dest := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
			h := headers.NewHeaders()

			hbreq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://httpbin.org/%s", dest), nil)
			if err != nil {
				log.Fatal("can't create request", "error", err)
			}
			res, err := clinet.Do(hbreq)
			defer res.Body.Close()
			if err != nil {
				err := w.WriteStatusLine(500)
				if err != nil {
					log.Fatal("error writing status line", "error", err)
				}
				h.Set("Connection", "close")
				return

			}
			err = w.WriteStatusLine(reponse.StatusCode(res.StatusCode))
			if err != nil {
				log.Fatal("cant write status line", "error", err)
			}
			for key, val := range res.Header {
				if key != "Content-Length" {
					h.Set(key, val[0])
				}
			}
			h.Set("Transfer-Encoding", "chunked")

			buf := make([]byte, 32)
			for {
				_, err := res.Body.Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						_, err := w.WriteChunkedBodyDone()
						if err != nil {
							log.Fatal("error writing chunked done", "error", err)
						}
						break
					}
					panic("we messed up response body reading")
				}
				_, err = w.WriteChunkedBody(buf)
				if err != nil {
					log.Fatal("error writing a chunk", "error", err)
				}
			}
			return
		}
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
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

		case "/myproblem":
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
		default:
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
