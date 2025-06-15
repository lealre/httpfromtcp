package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// Reads 8 bytes at a time
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for {

		buffer := make([]byte, 8)
		bytesRead, err := file.Read(buffer)

		if err != nil {
			if err == io.EOF {
				// log.Print("EOF reached...")
				break
			}
			log.Fatal(err)
		}

		text := string(buffer[:bytesRead])
		fmt.Printf("read: %s\n", text)
	}

}
