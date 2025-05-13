# streamcrypt

A fast and composable Go library for streaming compression and encryption. Designed for large file
processing pipelines, backups, and secure data storage or transfer.

[![License](https://img.shields.io/github/license/hashmap-kz/streamcrypt)](https://github.com/hashmap-kz/streamcrypt/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hashmap-kz/streamcrypt)](https://goreportcard.com/report/github.com/hashmap-kz/streamcrypt)
[![Workflow Status](https://img.shields.io/github/actions/workflow/status/hashmap-kz/streamcrypt/ci.yml?branch=master)](https://github.com/hashmap-kz/streamcrypt/actions/workflows/ci.yml?query=branch:master)
[![GitHub Issues](https://img.shields.io/github/issues/hashmap-kz/streamcrypt)](https://github.com/hashmap-kz/streamcrypt/issues)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hashmap-kz/streamcrypt)](https://github.com/hashmap-kz/streamcrypt/blob/master/go.mod#L3)
[![Latest Release](https://img.shields.io/github/v/release/hashmap-kz/streamcrypt)](https://github.com/hashmap-kz/streamcrypt/releases/latest)

---

## ✨ Features

- ✅ Stream-based compression and encryption (no full reads into memory)
- ✅ Pluggable compressors (`gzip`, `zstd`)
- ✅ Pluggable encryption backends (default: `AES-256-GCM` with Argon2 key derivation)
- ✅ Chunked encryption: safer and faster on large inputs
- ✅ Clean, testable design with `io.Reader/io.Writer` pipelines

### 💻 Usage

You can use `streamcrypt` directly in your Go code as a streaming compression/encryption library:

```
// Having a typical storage intefrace (impl may be any: s3, sftp, etc...)

type Storage interface {
	PutObject(ctx context.Context, path string, r io.Reader) error
	ReadObject(ctx context.Context, path string) (io.ReadCloser, error)
}

// Having a 'repo', that wraps conpression/encryption on streams: 

type repoImpl struct {
	storage    storage.Storage  // required: e.g. LocalImpl()
	compressor codec.Compressor // optional
	crypter    crypt.Crypter    // optional
}

func (repo *repoImpl) PutObject(ctx context.Context, path string, r io.Reader) (string, error) {
	var err error
	fullPath := repo.encodePath(path)

	// Compress and encrypt
	encReader, err := pipe.CompressAndEncryptOptional(r, repo.compressor, repo.crypter)
	if err != nil {
		return "", err
	}

	// Store in repo
	err = repo.storage.PutObject(ctx, fullPath, encReader)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

func (repo *repoImpl) ReadObject(ctx context.Context, path string) (io.ReadCloser, error) {
	var err error
	fullPath := repo.encodePath(path)

	// Open() that needs to be closed
	obj, err := repo.storage.ReadObject(ctx, fullPath)
	if err != nil {
		return nil, err
	}

	var dec codec.Decompressor
	if repo.compressor != nil {
		dec = codec.GetDecompressor(repo.compressor)
		if dec == nil {
			obj.Close()
			return nil, fmt.Errorf("cannot decide decompressor for: %s", repo.compressor.FileExtension())
		}
	}

	readCloser, err := pipe.DecryptAndDecompressOptional(obj, repo.crypter, dec)
	if err != nil {
		obj.Close()
		return nil, err
	}

	return ioutils.NewMultiCloser(readCloser, obj, readCloser), nil
}
```

---

## 🧩 Project Structure

| Package   | Purpose                                    |
|-----------|--------------------------------------------|
| `codec/`  | Pluggable compressors (gzip, etc.)         |
| `crypt/`  | Pluggable encryption implementations       |
| `pipe/`   | The core streaming pipeline                |
| `aesgcm/` | Chunked AES-GCM with Argon2 key derivation |

---

## 🔐 Security

- Uses **AES-256-GCM** for authenticated encryption
- Keys are derived via **Argon2id** with a random salt
- Each chunk is encrypted independently with unique nonce

---

## 📄 License

MIT License. See [LICENSE](./LICENSE) for details.

---

## 🙌 Acknowledgements

- [`filippo.io/age`](https://pkg.go.dev/filippo.io/age) – for inspiration
