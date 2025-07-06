package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/lealre/httpfromtcp/internal/headers"
	"github.com/lealre/httpfromtcp/internal/request"
	"github.com/lealre/httpfromtcp/internal/response"
	"github.com/lealre/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, testHandler)
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

func testHandler(w *response.Writer, req *request.Request) {
	html400 := `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

	if req.RequestLine.RequestTarget == "/yourproblem" {
		// resp line
		w.WriteStatusLine(response.BadRequest)

		// headers
		headers := headers.NewHeaders()
		contentSize := strconv.Itoa(len(html400))
		headers.Set("content-length", contentSize)
		headers.Set("content-type", "text/html")
		w.WriteHeaders(headers)

		// body
		w.WriteBody([]byte(html400))

	}

	html500 := `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

	if req.RequestLine.RequestTarget == "/myproblem" {
		// resp line
		w.WriteStatusLine(response.InternalServerError)

		// headers
		headers := headers.NewHeaders()
		contentSize := strconv.Itoa(len(html500))
		headers.Set("content-length", contentSize)
		headers.Set("content-type", "text/html")
		w.WriteHeaders(headers)

		// body
		w.WriteBody([]byte(html500))
	}

	html200 := `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

	// resp line
	w.WriteStatusLine(response.Ok)

	// headers
	headers := headers.NewHeaders()
	contentSize := strconv.Itoa(len(html200))
	headers.Set("content-length", contentSize)
	headers.Set("content-type", "text/html")
	w.WriteHeaders(headers)

	// body
	w.WriteBody([]byte(html200))
}
