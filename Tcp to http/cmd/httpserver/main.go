package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	server "tcpTohttp/internal"
	"tcpTohttp/internal/headers"
	"tcpTohttp/internal/request"
	"tcpTohttp/internal/response"
)

const port = 42069

func handlerfunc(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		data := []byte(`
		<html>
			<head>
				<title>400 Bad Request</title>
			</head>
			<body>
				<h1>Bad Request</h1>
				<p>Your request honestly kinda sucked.</p>
			</body>
		</html>`)
		headers := headers.NewHeaders()
		headers.Set("Content-Type", "text/html")
		headers.Set("Content-Length", strconv.Itoa(len(data)))
		w.WriteStatusLine(400)
		w.WriteHeaders(headers)
		w.WriteBody(data)

	case "/myproblem":
		data := []byte(`
		<html>
			<head>
				<title>500 Internal Server Error</title>
			</head>
			<body>
				<h1>Internal Server Error</h1>
    			<p>Okay, you know what? This one is on me.</p>
			</body>
		</html>`)
		headers := headers.NewHeaders()
		headers.Set("Content-Length", strconv.Itoa(len(data)))
		headers.Set("Content-Type", "text/html")
		w.WriteStatusLine(500)
		w.WriteHeaders(headers)
		w.WriteBody(data)

	default:
		data := []byte(`
		<html>
			<head>
				<title>200 OK</title>
			</head>
			<body>
				<h1>Success!</h1>
    			<p>Your request was an absolute banger.</p>
			</body>
		</html>`)
		headers := headers.NewHeaders()
		headers.Set("Content-Length", strconv.Itoa(len(data)))
		headers.Set("Content-Type", "text/html")
		w.WriteStatusLine(200)
		w.WriteHeaders(headers)
		w.WriteBody(data)
	}

}
func main() {
	server, err := server.Serve(port, handlerfunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
