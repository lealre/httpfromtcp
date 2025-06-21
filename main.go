package main

import (
	"io"
	"log"
	"strings"
)

func main() {
	// Reads 8 bytes at a time
	// file, err := os.Open("messages.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// linesChn := getLinesChannel(file)

	// for line := range linesChn {
	// 	fmt.Printf("read: %s\n", line)
	// }

	tcpReader()

	// compareRead()

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	var currentLine string
	ch := make(chan string, 100)
	go func() {
		defer close(ch)

		for {
			buffer := make([]byte, 8)
			bytesRead, err := f.Read(buffer)

			if err != nil {
				if err == io.EOF {
					if currentLine != "" {
						ch <- currentLine
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
				ch <- currentLine
				currentLine = subStrings[1]
				continue
			}
			currentLine += text

		}
	}()

	return ch
}
