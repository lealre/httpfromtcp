package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func compareRead() {
	filePath := "./test.pdf"

	// Method 1: Read entire file at once (for comparison)
	start := time.Now()
	readEntireFile(filePath)
	fmt.Printf("\nRead entire file at once: %v\n", time.Since(start))

	// Method 2: Read in 1KB chunks
	start = time.Now()
	readInChunks(filePath, 1024) // 1KB buffer
	fmt.Printf("Read in 1KB chunks: %v\n", time.Since(start))
}

// Read entire file (simplest, but risky for large files)
func readEntireFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read %d bytes (all at once)\n", len(data))
}

// Read file in chunks (safe for any size)
func readInChunks(path string, bufferSize int) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buffer := make([]byte, bufferSize)
	totalBytes := 0

	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		totalBytes += bytesRead
		// Process buffer[:bytesRead] here if needed
	}
	fmt.Printf("Read %d bytes (in %d-byte chunks)\n", totalBytes, bufferSize)
}
