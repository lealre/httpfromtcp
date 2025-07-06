package response

import (
	"fmt"
	"io"

	"github.com/lealre/httpfromtcp/internal/headers"
)

type Writer struct {
	Writer io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	WriteStatusLine(w.Writer, statusCode)
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for key, value := range headers {
		keyPairHeaderValue := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Writer.Write([]byte(keyPairHeaderValue))
		if err != nil {
			return err
		}
	}
	w.Writer.Write([]byte("\r\n"))
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}
