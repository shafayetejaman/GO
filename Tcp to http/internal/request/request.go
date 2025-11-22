package request

import (
	"bytes"
	"errors"
	"io"
	"log"
	"log/slog"
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

func (r *Request) parse(data []byte) (int, error) {
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
			currDataLen := len(data)
			slog.Info("statebody", "data", data, "bodyLen", r.BodyLen)

			if r.BodyLen == -1 && currDataLen > 0 {
				return 0, errors.New("content length not found")
			}

			if currDataLen == 0 {
				if r.BodyLen > 0 {
					return 0, errors.New("missing body")

				}
				r.State = StateDone
				return 0, nil
			}

			if r.BodyLen == 0 && currDataLen > 0 {
				return 0, errors.New("the body is too big")

			}
			read := min(r.BodyLen, currDataLen)
			r.Body = append(r.Body, data[:read]...)
			r.BodyLen -= read

			return read, nil

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
	isEOF := false

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
					isEOF = true
				} else {
					return nil, err
				}
			}
		}
		read += n
		slog.Info("readFom", "data", buffer[:read], "iseof", isEOF)
		n, err = req.parse(buffer[:read])
		if err != nil {
			return nil, err
		}
		copy(buffer, buffer[n:read])
		read -= n

	}

	return &req, nil
}
