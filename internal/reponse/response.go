package reponse

import (
	"fmt"
	"http-serveur/internal/headers"
	"io"
	"strconv"
)

type StatusCode = uint16

const (
	StatusOk          StatusCode = 200
	StatusBadRequest  StatusCode = 400
	StatusServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOk:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
	case StatusServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
	default:
		_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)))
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	iContentLength := strconv.Itoa(contentLen)
	h.Set("Content-Length", iContentLength)
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return *h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headerLine := ""
	headers.ForEach(func(n, v string) {
		headerLine += fmt.Sprintf("%s: %s", n, v)
		headerLine += "\r\n"
	})
	headerLine += "\r\n\r\n"
	w.Write([]byte(headerLine))
	return nil
}
