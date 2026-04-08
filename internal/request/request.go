package request

import (
	"bytes"
	"fmt"
	"http-serveur/internal/headers"
	"io"
	"slices"
	"strconv"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
	Headers     *headers.Headers
	Body        []byte
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    []byte(""),
	}
}

var ERROR_BAD_REQUEST_LINE = fmt.Errorf("Bad Request Line")
var SEPERATOR = []byte("\r\n")
var ERROR_BAD_HTTP_VERSION = fmt.Errorf("Bad Http Version")

func (r *Request) hasBody() bool {
	lengthStr, ok := r.Headers.Get("content-length")
	if !ok {
		return false
	}
	return lengthStr != "0"
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		if len(data[read:]) == 0 {
			break outer
		}
		switch r.state {

		case StateInit:
			requestLine, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *requestLine
			read += n

			r.state = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n

			if done {
				if r.hasBody() {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}

		case StateBody:
			bodyLengthStr, ok := r.Headers.Get("content-length")
			if !ok {
				panic("chunked not implemented")
			}
			bodyLength, err := strconv.Atoi(bodyLengthStr)
			if err != nil {
				r.state = StateDone
			}
			if bodyLength == 0 {
				panic("chunked not implemented")
			}

			remaining := min(bodyLength-len(r.Body), len(data[read:]))
			r.Body = append(r.Body, data[read:read+remaining]...)
			read += remaining
			if len(r.Body) == bodyLength {
				r.state = StateDone
			}

		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	startLine := b[:idx]
	read := idx + len(SEPERATOR)
	parts := bytes.Split(startLine, []byte{' '})
	if len(parts) != 3 {
		return nil, 0, ERROR_BAD_REQUEST_LINE
	}
	versionParts := bytes.Split(parts[2], []byte{'/'})
	if len(versionParts) != 2 || !slices.Equal(versionParts[0], []byte("HTTP")) || !slices.Equal(versionParts[1], []byte("1.1")) {
		return nil, 0, ERROR_BAD_REQUEST_LINE
	}
	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(versionParts[1]),
	}

	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buff := make([]byte, 1024)
	buffLength := 0
	for !request.done() {
		n, err := reader.Read(buff[buffLength:])
		if err != nil {
			return nil, err
		}

		buffLength += n
		readN, err := request.parse(buff[:buffLength])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[readN:buffLength])
		buffLength -= readN

	}
	return request, nil
}
