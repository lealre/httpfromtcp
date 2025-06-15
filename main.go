package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// Reads 8 bytes at a time
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var currentLine string

	for {

		buffer := make([]byte, 8)
		bytesRead, err := file.Read(buffer)

		if err != nil {
			if err == io.EOF {
				if currentLine != "" {
					fmt.Printf("read: %s\n", currentLine)
				}
				// log.Print("EOF reached...")
				break
			}
			log.Fatal(err)
		}

		text := string(buffer[:bytesRead])

		if strings.Contains(text, "\n") {
			subStrings := strings.Split(text, "\n")
			currentLine += subStrings[0]
			fmt.Printf("read: %s\n", currentLine)
			currentLine = subStrings[1]
			continue
		}
		currentLine += text

	}

}
