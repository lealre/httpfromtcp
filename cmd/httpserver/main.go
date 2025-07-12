package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lealre/httpfromtcp/internal/request"
	"github.com/lealre/httpfromtcp/internal/response"
	"github.com/lealre/httpfromtcp/internal/server"
)

const port = 42069
const bufferSize = 1024

func main() {
	server, err := server.Serve(port, handler)
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

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		handlerChunkEncoding(w, req)
		return
	}
	handler200(w, req)
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.BadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.InternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.Ok)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handlerChunkEncoding(w *response.Writer, req *request.Request) {
	path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := fmt.Sprintf("https://httpbin.org/%s", path)
	resp, err := http.Get(url)
	if err != nil {
		errorBody := fmt.Sprintf("error executing endpoint: %s", err)
		w.WriteStatusLine(response.InternalServerError)
		h := response.GetDefaultHeaders(len(errorBody))
		w.WriteHeaders(h)
		w.WriteBody([]byte(errorBody))
		return
	}

	defer resp.Body.Close()

	w.WriteStatusLine(response.Ok)
	h := response.GetDefaultHeaders(0)
	h.Remove("content-length")
	h.Remove("connection")
	h.Override("Content-Type", "application/json")
	h.Set("Transfer-Encoding", "chunked")
	w.WriteHeaders(h)

	buff := make([]byte, bufferSize)
	for {
		n, err := resp.Body.Read(buff)
		fmt.Println("Read", n, "bytes")
		if n > 0 {
			_, err = w.WriteChunkedBody(buff[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("Error writing chunked body done:", err)
	}
}
