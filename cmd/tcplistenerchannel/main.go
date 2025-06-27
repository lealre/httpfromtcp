package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

/*
TCP Line Reader Server

Usage:

 1. Start server (with output logging):
    go run ./cmd/tcplineserver | tee /tmp/tcp.txt

 2. Send test messages:
    printf "Hello\nWorld\n" | nc localhost 42069
    -or-
    echo "message" | nc localhost 42069

 3. View logs:
    cat /tmp/tcp.txt

Behavior:
- Listens on TCP port 42069
- Accepts multiple connections sequentially (not concurrently)
- Reads until connection closes (EOF)
- Buffers all received data before sending complete message through channel
- Prints:
  - Connection events (log.Printf)
  - Raw message content (fmt.Println)

- Closes connection after EOF
*/
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

		// Read and process the connection
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}

		conn.Close()
		log.Printf("Connection from %s closed", conn.RemoteAddr().String())
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string, 100)
	go func() {
		defer close(ch)
		defer f.Close()

		buf := make([]byte, 1024)
		var buffer strings.Builder

		for {
			n, err := f.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("Read error: %v", err)
				}
				break
			}

			buffer.Write(buf[:n])
		}

		// Send the complete received message
		ch <- buffer.String()
	}()
	return ch
}
