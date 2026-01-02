package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
	Headers     map[string]string
	Body        []byte
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

var ERROR_BAD_REQUEST_LINE = fmt.Errorf("Bad Request Line")
var SEPERATOR = "\r\n"
var ERROR_BAD_HTTP_VERSION = fmt.Errorf("Bad Http Version")

func (r *Request) parse(data []byte) (int, error) {

}
func (r *Request) done() bool {
	return r.state == StateDone
}

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, b, nil
	}
	startLine := b[:idx]
	restOfMessage := b[idx+len(SEPERATOR):]
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMessage, ERROR_BAD_REQUEST_LINE
	}
	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 || versionParts[0] != "HTTP" || versionParts[1] != "1.1" {
		return nil, restOfMessage, ERROR_BAD_REQUEST_LINE
	}
	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   versionParts[1],
	}

	return rl, restOfMessage, nil
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
