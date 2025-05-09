package boot

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	original := []byte("this is some secret data that should roundtrip")
	password := "s3cr3t"

	// Encrypt
	encReader, err := Encrypt(bytes.NewReader(original), password)
	require.NoError(t, err)

	// Decrypt
	decReader, err := Decrypt(encReader, password)
	require.NoError(t, err)
	defer decReader.Close()

	result, err := io.ReadAll(decReader)
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestEncrypt_EmptyInput(t *testing.T) {
	password := "emptytest"
	var buf bytes.Buffer

	encReader, err := Encrypt(&buf, password)
	require.NoError(t, err)

	decReader, err := Decrypt(encReader, password)
	require.NoError(t, err)
	defer decReader.Close()

	result, err := io.ReadAll(decReader)
	require.NoError(t, err)
	assert.Len(t, result, 0)
}
