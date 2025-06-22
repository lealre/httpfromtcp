package request

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	var request Request
	message := string(bytes)
	firstLine := strings.Split(message, "\r\n")[0]

	requestLine, err := parserRequestLine(firstLine)
	if err != nil {
		return &Request{}, err
	}
	request.RequestLine = requestLine
	return &request, nil

}

func parserRequestLine(s string) (RequestLine, error) {

	var requestLine RequestLine
	fmt.Printf("String %s\n", s)
	// method
	components := strings.Split(s, " ")
	fmt.Printf("Component\n")
	fmt.Println(components)

	method := components[0]
	if method != "GET" && method != "POST" {
		return RequestLine{}, errors.New("invalid Method")
	}
	requestLine.Method = method

	// target
	target := components[1]
	if !strings.HasPrefix(target, "/") {
		return RequestLine{}, errors.New("invalid target")
	}
	requestLine.RequestTarget = target

	// version
	version := components[2]
	if version != "HTTP/1.1" {
		return RequestLine{}, errors.New("invalid HTTP version")
	}
	requestLine.HttpVersion = strings.Split(version, "/")[1]

	return requestLine, nil
}
