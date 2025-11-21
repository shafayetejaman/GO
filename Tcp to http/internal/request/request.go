package request

import (
	"bytes"
	"errors"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"
	"tcpTohttp/internal/headers"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       State
	Body        []byte
	BodyLen     int
}

type State int

const (
	StateInitialized State = iota
	StateDone
	StateParsingHeaders
	StateParseBody
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const CRLF = "\r\n"

func (r *Request) parse(data []byte, eof bool) (int, error) {
	for {
		switch r.State {
		case StateDone:
			return 0, nil

		case StateInitialized:
			reqLine, n, err := parseRequestLine(data)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				return 0, nil
			}
			r.RequestLine = reqLine
			r.State = StateParsingHeaders
			return n, nil

		case StateParsingHeaders:
			// slog.Info("stateParheer#66", "data", data)

			n, done, err := r.Headers.Parse(data)

			if err != nil {
				return 0, err
			}
			if !done {
				return n, nil
			}
			r.State = StateParseBody
			contentLen := r.Headers.Get("Content-Length")

			if contentLen != "" {
				n, err := strconv.Atoi(contentLen)
				if err != nil {
					return 0, err
				}
				r.BodyLen = n
			}

			return n, nil

		case StateParseBody:
			currBodyLen := len(data)

			if r.BodyLen == -1 && currBodyLen > 0 {
				return 0, errors.New("the body is too big")
			}

			if r.BodyLen == -1 || (r.BodyLen == 0 && currBodyLen == 0) {
				r.State = StateDone
				return 0, nil
			}

			if r.BodyLen > 0 && len(data) <= 0 && eof {

				return 0, errors.New("missing body")
			}

			r.Body = append(r.Body, data...)
			r.BodyLen -= len(data)

			return len(data), nil

		default:
			log.Fatal("state dose not match")
		}

	}
}

func (r Request) done() bool {
	return r.State == StateDone
}

func parseRequestLine(data []byte) (RequestLine, int, error) {

	var reqLine []byte
	req := RequestLine{}
	read := 0
	if i := bytes.Index(data, []byte(CRLF)); i == -1 {
		return req, 0, nil
	} else {
		reqLine = data[:i]
		read = i + len(CRLF)
	}

	reqPart := strings.Split(string(reqLine), " ")

	if len(reqPart) != 3 {
		return req, 0, errors.New("missing or to mangy arguments")
	}
	// GET /coffee HTTP/1.1
	// validate method
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	if !slices.Contains(methods, reqPart[0]) {
		return req, 0, errors.New("invalid method")
	}
	req.Method = reqPart[0]

	// validate http vertion
	httpv := strings.Split(reqPart[2], "/")
	if !(len(httpv) == 2 && httpv[0] == "HTTP" && httpv[1] == "1.1") {
		return req, 0, errors.New("invalid http vertions")
	}
	req.HttpVersion = httpv[1]

	// validate path
	path := reqPart[1]
	invalidChar := []string{"\\", " ", "+", "\n", "<", ">",
		"|", "\"", "\\'", "{", "}", "^"}

	if !strings.HasPrefix(path, "/") ||
		slices.Contains(invalidChar, path) ||
		!isAssci(path) {
		return req, 0, errors.New("invalid path")
	}
	req.RequestTarget = path

	return req, read, nil
}
func isAssci(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}
func RequestFromReader(reader io.Reader) (*Request, error) {
	req := Request{State: StateInitialized,
		Headers: *headers.NewHeaders(),
		BodyLen: -1}
	buffer := make([]byte, 8)
	read := 0
	var isEOF bool

	for !req.done() {
		var n int
		var err error
		if read == len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
		if !isEOF {
			n, err = reader.Read(buffer[read:])
			if err != nil {
				if err == io.EOF {
					println(err)
					isEOF = true
				} else {
					return nil, err
				}
			}
		}
		read += n
		n, err = req.parse(buffer[:read], isEOF)
		if err != nil {
			return nil, err
		}
		copy(buffer, buffer[n:read])
		read -= n

	}

	return &req, nil
}
