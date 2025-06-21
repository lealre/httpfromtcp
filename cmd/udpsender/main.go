package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

/*
UDP Sender

Usage:

 1. First terminal (receiver):
    nc -u -l 42069                # Basic UDP listener
    -or-
    socat -u UDP-RECV:42069 -     # Better formatted output

 2. Second terminal (sender):
    go run ./cmd/udpsender        # Starts interactive prompt

Behavior:
- Reads input line-by-line from stdin (including newlines)
- Sends each line as a separate UDP packet to localhost:42069
- Preserves all whitespace and newline characters
- Runs until interrupted with Ctrl+C

Note: For testing binary data, use:

	nc -u -l 42069 | hexdump -C
*/
func main() {
	// Resolve UDP address
	udpAddress, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create reader for stdin
	reader := bufio.NewReader(os.Stdin)

	// Main loop
	for {
		// Print prompt
		print("> ")

		// Read input
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}

		// Trim newline and check for empty input
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// Send over UDP
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}
