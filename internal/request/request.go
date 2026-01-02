package request

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

func (r *RequestLine) ValidHttp() bool {
	return r.HttpVersion == "HTML/1.1"
}

var ERROR_BAD_REQUEST_LINE = fmt.Errorf("Bad Request Line")
var SEPERATOR = "\r\n"
var ERROR_BAD_HTTP_VERSION = fmt.Errorf("Bad Http Version")

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
	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}
	if !rl.ValidHttp() {
		return nil, restOfMessage, ERROR_BAD_HTTP_VERSION
	}
	return rl, restOfMessage, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("couldnt read request")
	}
}
