package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func tcpReader() {
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

func getLinesChannel_(f io.ReadCloser) <-chan string {
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
