# Certificate, HTTPS & HTTP/2

## Overview

| Istilah | Arti |
|---------|------|
| **HTTPS** | HTTP + TLS. Data terenkripsi antara client dan server. |
| **TLS** | Protokol keamanan yang mengenkripsi koneksi. |
| **Certificate** | File publik (cert.pem) — identitas server. |
| **Private Key** | File rahasia (key.pem) — hanya server yang punya. |
| **HTTP/2** | Versi HTTP yang lebih cepat. **WAJIB** pakai HTTPS (browser hanya support HTTP/2 via TLS). |

## Generate Self-Signed Certificate

Untuk development, kita bikin certificate sendiri (self-signed), bukan dari CA resmi.

```bash
openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 365
```

Penjelasan tiap flag:

| Flag | Arti |
|------|------|
| `req` | Certificate Signing Request (CSR) |
| `-x509` | Langsung output certificate (skip CSR) |
| `-newkey rsa:2048` | Generate private key baru, RSA 2048 bit |
| `-nodes` | No DES — private key **tanpa** password |
| `-keyout key.pem` | Simpan private key ke file `key.pem` |
| `-out cert.pem` | Simpan certificate ke file `cert.pem` |
| `-days 365` | Masa berlaku 1 tahun |

File yang dihasilkan:

```
├── cert.pem   ✅ boleh di-commit (isi publik)
├── key.pem    ❌ JANGAN di-commit (rahasia!)
```

## Cara Kerja TLS di Go

```
Client request ──→ Server punya cert.pem + key.pem
                       │
                       ▼
Server kirim cert.pem ke client ──→ Client verifikasi
                       │
                       ▼
Client + Server bikin session key (encrypted handshake)
                       │
                       ▼
Semua data terenkripsi dari sini
```

## Kode Go untuk HTTPS + HTTP/2

### 1. Import Package

```go
import (
    "crypto/tls"                // TLS configuration

    "golang.org/x/net/http2"    // enable HTTP/2 protocol
)
```

### 2. Load Certificate

```go
cert := "cert.pem"   // file certificate (self-signed)
key := "key.pem"     // file private key (jangan di-share)
```

### 3. Konfigurasi TLS

```go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,  // minimal TLS 1.2 — tolak versi lama (1.0, 1.1)
}
```

### 4. Custom HTTP Server dengan TLS

```go
server := &http.Server{
    Addr:      ":3000",
    Handler:   nil,           // nil = pake DefaultServeMux
    TLSConfig: tlsConfig,     // pasang config TLS
}
```

### 5. Enable HTTP/2

```go
http2.ConfigureServer(server, &http2.Server{})
```

Baris ini mengaktifkan HTTP/2 di atas TLS. Tanpa ini, server cuma HTTP/1.1.

### 6. Start HTTPS Server

```go
err := server.ListenAndServeTLS(cert, key)
if err != nil {
    log.Fatalln("could not start server", err)
}
```

`ListenAndServeTLS` mirip `ListenAndServe` tapi dengan TLS:
- Parameter 1: path certificate (`cert.pem`)
- Parameter 2: path private key (`key.pem`)

## Testing

### Test dengan curl

```bash
# HTTPS (tanpa verifikasi cert — karena self-signed)
curl -k https://localhost:3000/orders

# HTTPS + HTTP/2
curl --http2 -k https://localhost:3000/orders -v
```

Flag `-k` = skip certificate verification (karena self-signed).

### Cek HTTP/2

Dari output `curl -v`, cari baris:

```
* Using HTTP/2
```

Kalau ada tulisan itu, HTTP/2 berjalan.

### Error yang Sering Muncul

| Error | Penyebab | Solusi |
|-------|----------|--------|
| `tls: first record does not look like a TLS handshake` | Client HTTP ke server HTTPS | Pakai `https://` bukan `http://` |
| `certificate is not trusted` | Self-signed, truststore tidak punya CA kita | Pakai `-k` di curl atau trust manual |
| `no such file or directory: cert.pem` | File cert tidak ada | Generate dulu dengan `openssl` |
