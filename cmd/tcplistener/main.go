package main

import (
	"fmt"
	"log"
	"net"

	"github.com/lealre/httpfromtcp/internal/request"
)

func main() {
	tcpListener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer tcpListener.Close()

	log.Printf("Server started, listening on %s", tcpListener.Addr().String())

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		log.Printf("New connection accepted from %s", conn.RemoteAddr().String())

		response, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", response.RequestLine.Method)
		fmt.Printf("- Target: %s\n", response.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", response.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range response.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Body:")
		fmt.Println(string(response.Body))

		conn.Close()
		log.Printf("Connection from %s closed", conn.RemoteAddr().String())
	}
}
