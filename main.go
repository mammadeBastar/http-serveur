package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannels(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)
		line := ""
		for {
			buff := make([]byte, 8)
			n, err := f.Read(buff)
			if err != nil {
				break
			}
			buff = buff[:n]
			if i := bytes.IndexByte(buff, '\n'); i != -1 {
				line += string(buff[:i])
				buff = buff[i+1:]
				out <- line
				line = ""
			}
			line += string(buff)
		}
		if len(line) != 0 {
			out <- line
		}
	}()

	return out
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		for line := range getLinesChannels(conn) {
			fmt.Printf("read: %s\n", line)
		}

	}
}
