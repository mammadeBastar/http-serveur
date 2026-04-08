package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid double header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nanotherMohsen:   mohsi   \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, _ := headers.Get("Host")
	assert.Equal(t, "localhost:42069", host)
	anotherMohsen, _ := headers.Get("anotherMohsen")
	assert.Equal(t, "mohsi", anotherMohsen)
	assert.Equal(t, 52, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid field name
	headers = NewHeaders()
	data = []byte("mحsen: oneMohsen\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: valid case insensetivity
	headers = NewHeaders()
	data = []byte("moHSenOne: mohsen\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	mohsenOne, _ := headers.Get("mohsenONE")
	assert.Equal(t, "mohsen", mohsenOne)
	assert.Equal(t, 21, n)
	assert.True(t, done)

	// Test: valid multiple headers with the same name
	headers = NewHeaders()
	data = []byte("moHSenOne: mohsen\r\nmohsenone: ahmad\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	mohsenOne, _ = headers.Get("mohsenONE")
	assert.Equal(t, "mohsen,ahmad", mohsenOne)
	assert.Equal(t, 39, n)
}
