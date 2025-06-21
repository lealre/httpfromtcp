package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

/*
File Reader Benchmark

Usage:

	go run ./cmd/filereader

Behavior:
- Reads test.pdf using two methods:
 1. Entire file at once (os.ReadFile)
 2. Chunked reading (1KB buffers)

- Prints timing results for comparison
- Shows total bytes read in each case

Key Differences:
+---------------+-----------------+-------------------+
| Method        | Memory Use      | Best For          |
+---------------+-----------------+-------------------+
| Entire file   | High (file size)| Small files       |
| 1KB chunks    | Low (1KB)       | Large files       |
+---------------+-----------------+-------------------+

Note: Change bufferSize in code to test different chunk sizes.
*/
func main() {
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
