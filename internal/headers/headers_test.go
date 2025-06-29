package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, len(data)-2, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	// First header (partial data)
	data1 := []byte("Host: localhost:42069\r\n")
	n, done, err = headers.Parse(data1)
	require.NoError(t, err)
	assert.Equal(t, len(data1), n)
	assert.False(t, done)
	assert.Equal(t, "localhost:42069", headers["host"])
	// Second header (remaining data)
	data2 := []byte("User-Agent: Go-HTTP-Parser\r\n\r\n")
	n, done, err = headers.Parse(data2)
	require.NoError(t, err)
	assert.Equal(t, len(data2)-2, n)
	assert.False(t, done)
	assert.Equal(t, "Go-HTTP-Parser", headers["user-agent"])

	// Test: Case-insensitive header merging
	headers = NewHeaders()
	// First header
	data1 = []byte("Set-Person: name1\r\n")
	n, done, err = headers.Parse(data1)
	require.NoError(t, err)
	assert.Equal(t, len(data1), n)
	assert.False(t, done)
	assert.Equal(t, "name1", headers["set-person"])
	// Second header
	data2 = []byte("set-person: name2\r\n\r\n")
	n, done, err = headers.Parse(data2)
	require.NoError(t, err)
	assert.Equal(t, len(data2)-2, n)
	assert.False(t, done)
	assert.Equal(t, "name1, name2", headers["set-person"])

}
