# streamcrypt

A fast and composable Go library and CLI tool for streaming compression and encryption. Designed for large file
processing pipelines, backups, and secure data storage or transfer.

---

## ✨ Features

- ✅ Stream-based compression and encryption (no full reads into memory)
- ✅ Pluggable compressors (`gzip` supported, `zstd` planned)
- ✅ Pluggable encryption backends (default: `AES-256-GCM` with Argon2 key derivation)
- ✅ Chunked encryption: safer and faster on large inputs
- ✅ Clean, testable design with `io.Reader/io.Writer` pipelines
- ✅ CLI support via `cobra`

---

## 📦 Usage

### 🔒 CLI: Encrypt a file

```bash
streamcrypt encrypt --in plain.txt --out secret.gz.aes --password "s3cr3t"
```

### 🔓 CLI: Decrypt a file

```bash
streamcrypt decrypt --in secret.gz.aes --out plain.txt --password "s3cr3t"
```

### 💡 Pipe data in/out (streaming)

```bash
cat plain.txt | streamcrypt encrypt --password "s3cr3t" > secret.gz.aes
cat secret.gz.aes | streamcrypt decrypt --password "s3cr3t" > plain.txt
```

---

## 🧩 Project Structure

| Package   | Purpose                                    |
|-----------|--------------------------------------------|
| `cmd/`    | Cobra-based CLI definitions                |
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
- [`spf13/cobra`](https://github.com/spf13/cobra) – CLI handling
