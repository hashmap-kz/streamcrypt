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
import (
    "bytes"
    "io"
    "log"

    "github.com/hashmap-kz/streamcrypt/pkg/boot"
)

func encryptAndDecryptExample() {
    input := []byte("stream me securely")
    password := "s3cr3t"

    // Encrypt
    encReader, err := boot.Encrypt(bytes.NewReader(input), password)
    if err != nil {
        log.Fatal("encryption failed:", err)
    }

    // Decrypt
    decReader, err := boot.Decrypt(encReader, password)
    if err != nil {
        log.Fatal("decryption failed:", err)
    }
    defer decReader.Close()

    output, err := io.ReadAll(decReader)
    if err != nil {
        log.Fatal("read failed:", err)
    }

    log.Printf("Decrypted content: %s", string(output))
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
