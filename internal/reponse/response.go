package reponse

import (
	"fmt"
	"http-serveur/internal/headers"
	"io"
	"strconv"
)

type StatusCode = uint16
type writerState = string

const (
	StateStatusLine writerState = "status line"
	StateHeaders    writerState = "headers"
	StateBody       writerState = "body"
	StateDone       writerState = "nothing"
)

type Writer struct {
	state writerState
	io.Writer
}

const (
	StatusOk          StatusCode = 200
	StatusBadRequest  StatusCode = 400
	StatusServerError StatusCode = 500
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		state:  StateStatusLine,
		Writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != StateStatusLine {
		return fmt.Errorf("you should be writing %s rn", w.state)
	}
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
	w.state = StateHeaders
	return nil
}

func GetDefaultHeaders(contentLen int, contentType string) headers.Headers {
	h := headers.NewHeaders()
	iContentLength := strconv.Itoa(contentLen)
	h.Set("Content-Length", iContentLength)
	h.Set("Connection", "close")
	h.Set("Content-Type", contentType)
	return *h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != StateHeaders {
		return fmt.Errorf("you should be writing %s rn", w.state)
	}
	headerLine := ""
	headers.ForEach(func(n, v string) {
		headerLine += fmt.Sprintf("%s: %s", n, v)
		headerLine += "\r\n"
	})
	headerLine += "\r\n\r\n"
	w.Write([]byte(headerLine))
	w.state = StateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != StateBody {
		return 0, fmt.Errorf("you should be writing %s rn", w.state)
	}
	w.state = StateDone
	return w.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkSize := []byte(fmt.Sprintf("%X\r\n", len(p)))
	var written int
	n1, err := w.Write(chunkSize)
	if err != nil {
		return 0, err
	}
	p = append(p, []byte("\r\n")...)
	n2, err := w.Write(p)
	if err != nil {
		return written, err
	}
	return n1 + n2, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.Write([]byte("0\r\n\r\n"))
}
