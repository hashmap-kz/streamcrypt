package aesgcm

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkedGCMCrypto_SaltReadFailure(t *testing.T) {
	password := "."
	crypter := &ChunkedGCMCrypter{Password: password}
	r := bytes.NewReader([]byte{}) // too short to contain salt
	decReader, err := crypter.Decrypt(r)

	assert.Error(t, err)
	assert.Nil(t, decReader)
}

func TestChunkedGCMCrypto_EncryptFunction_HeaderAndSalt(t *testing.T) {
	password := "header-test"
	crypto := &ChunkedGCMCrypter{Password: password}
	var out bytes.Buffer

	writer, err := crypto.Encrypt(&out)
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	// Write and close with dummy data
	_, err = writer.Write([]byte("abc"))
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	written := out.Bytes()
	assert.True(t, bytes.HasPrefix(written, []byte("AEADv1")), "header prefix missing")
	assert.True(t, len(written) > len("AEADv1")+saltSize, "not enough data written")
}

func TestChunkedGCMCrypto_EncryptWriteFlushCloseBehavior(t *testing.T) {
	password := "flush-check"
	crypto := &ChunkedGCMCrypter{Password: password}
	var out bytes.Buffer

	writer, err := crypto.Encrypt(&out)
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	// Write data less than chunk size
	sample := bytes.Repeat([]byte("X"), 100)
	_, err = writer.Write(sample)
	assert.NoError(t, err)

	// Flush remaining data
	err = writer.Close()
	assert.NoError(t, err)
	assert.Greater(t, len(out.Bytes()), 0)
}

func TestChunkedGCMCrypto_DecryptFunction_InvalidHeader(t *testing.T) {
	password := "bad-header"
	crypto := &ChunkedGCMCrypter{Password: password}
	data := append([]byte("BADHDR"), make([]byte, 100)...) // malformed header
	reader, err := crypto.Decrypt(bytes.NewReader(data))
	assert.Nil(t, reader)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid file header")
}

func TestChunkedGCMCrypto_ChunkedWriterFlushChunkBoundary(t *testing.T) {
	password := "boundary-check"
	crypto := &ChunkedGCMCrypter{Password: password}
	var out bytes.Buffer

	writer, err := crypto.Encrypt(&out)
	assert.NoError(t, err)

	data := bytes.Repeat([]byte("A"), chunkSize)
	n, err := writer.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.NoError(t, writer.Close())
}

func TestChunkedGCMCrypto_ChunkedWriterMultipleFlushes(t *testing.T) {
	password := "multi-flush"
	crypto := &ChunkedGCMCrypter{Password: password}
	var out bytes.Buffer

	writer, err := crypto.Encrypt(&out)
	assert.NoError(t, err)

	chunk := bytes.Repeat([]byte("Z"), chunkSize)
	for i := 0; i < 3; i++ {
		_, err := writer.Write(chunk)
		assert.NoError(t, err)
	}
	assert.NoError(t, writer.Close())
	assert.Greater(t, len(out.Bytes()), 3*chunkSize) // includes nonce + tag overheads
}

func TestChunkedGCMCrypto_EncryptDecryptMultipleRandomFiles(t *testing.T) {
	tmpDir := t.TempDir()
	numFiles := 50
	maxSize := 512 * 1024 // 512 KB max

	password := "test-password"
	crypter := &ChunkedGCMCrypter{Password: password}

	for i := 0; i < numFiles; i++ {
		// Random file size (including zero)
		size := randomInt(0, maxSize)
		original := make([]byte, size)
		_, err := rand.Read(original)
		assert.NoError(t, err)

		// Compute original hash
		originalHash := sha256.Sum256(original)

		// Encrypt to buffer
		var encrypted bytes.Buffer
		encWriter, err := crypter.Encrypt(&encrypted)
		assert.NoError(t, err)
		_, err = encWriter.Write(original)
		assert.NoError(t, err)
		assert.NoError(t, encWriter.Close())

		// Decrypt from buffer
		decReader, err := crypter.Decrypt(bytes.NewReader(encrypted.Bytes()))
		assert.NoError(t, err)
		decrypted, err := io.ReadAll(decReader)
		assert.NoError(t, err)

		// Compute hash of decrypted
		decryptedHash := sha256.Sum256(decrypted)

		// Assert hash matches
		assert.Equal(t, originalHash, decryptedHash, "file %d: hash mismatch", i)

		// Optionally write a debug copy if the test fails
		if !assert.Equal(t, original, decrypted) {
			//nolint:errcheck
			_ = os.WriteFile(filepath.Join(tmpDir, "fail-original.bin"), original, 0o600)
			//nolint:errcheck
			_ = os.WriteFile(filepath.Join(tmpDir, "fail-decrypted.bin"), decrypted, 0o600)
			t.Fatalf("mismatch in file %d", i)
		}
	}
}

func randomInt(xmin, xmax int) int {
	if xmax <= xmin {
		return xmin
	}
	b := make([]byte, 4)
	//nolint:errcheck
	_, _ = rand.Read(b)
	return xmin + int(b[0])%(xmax-xmin)
}
