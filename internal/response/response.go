package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/lealre/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case 200:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case 400:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case 500:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		return nil
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("content-length", strconv.Itoa(contentLen))
	headers.Set("connection", "close")
	headers.Set("content-type", "text/plain")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		keyPair := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Write([]byte(keyPair))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}
