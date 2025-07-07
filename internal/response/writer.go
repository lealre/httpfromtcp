package response

import (
	"fmt"
	"io"

	"github.com/lealre/httpfromtcp/internal/headers"
)

type responseWriterStatus int

const (
	writerStarted responseWriterStatus = iota
	statusLineDone
	headersDone
	bodyDone
)

type Writer struct {
	Writer       io.Writer
	writerStatus responseWriterStatus
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer:       w,
		writerStatus: writerStarted,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerStatus != writerStarted {
		return fmt.Errorf("trying to write the reponse in the wrong order")
	}

	defer func() { w.writerStatus = statusLineDone }()

	WriteStatusLine(w.Writer, statusCode)
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerStatus != statusLineDone {
		return fmt.Errorf("trying to write the reponse in the wrong order")
	}

	defer func() { w.writerStatus = headersDone }()

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
	if w.writerStatus != headersDone {
		return 0, fmt.Errorf("trying to write the reponse in the wrong order")
	}

	defer func() { w.writerStatus = bodyDone }()
	n, err := w.Writer.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}
