package boot

import (
	"io"

	"github.com/hashmap-kz/streamcrypt/v1/pkg/codec"
	"github.com/hashmap-kz/streamcrypt/v1/pkg/crypt/aesgcm"
	"github.com/hashmap-kz/streamcrypt/v1/pkg/pipe"
)

func Encrypt(in io.Reader, password string) (io.Reader, error) {
	// Compression setup
	compressor := codec.GzipCompressor{}

	// Encryption setup
	crypter := aesgcm.NewChunkedGCMCrypter(password)

	return pipe.CompressAndEncryptOptional(in, compressor, crypter)
}

func Decrypt(in io.Reader, password string) (io.ReadCloser, error) {
	// Decompression setup
	decompressor := codec.GzipDecompressor{}

	// Decryption setup
	crypter := aesgcm.NewChunkedGCMCrypter(password)

	return pipe.DecryptAndDecompressOptional(in, crypter, decompressor)
}
