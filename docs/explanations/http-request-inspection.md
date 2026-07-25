# HTTP Request Inspection di Go

## `r.Proto` — HTTP Version

`r.Proto` adalah string yang menandakan HTTP version dari request client.

```go
httpVersion := r.Proto
fmt.Println(httpVersion)
```

| Nilai `r.Proto` | Arti                          |
| --------------- | ----------------------------- |
| `"HTTP/1.1"`    | Client pake HTTP/1.1          |
| `"HTTP/2"`      | Client pake HTTP/2            |
| `"HTTP/1.0"`    | Client pake HTTP/1.0 (jarang) |

HTTP/2 request dikirim sebagai `HTTP/2`, bukan `HTTP/2.0` — beda dengan HTTP/1.x.

---

## `r.TLS` — TLS Connection Info

`r.TLS` adalah pointer ke `tls.ConnectionState`. Nilai:

```go
if r.TLS != nil {
    // koneksi HTTPS (TLS aktif)
    fmt.Println("TLS version:", r.TLS.Version)
    fmt.Println("Cipher suite:", r.TLS.CipherSuite)
    fmt.Println("Server name:", r.TLS.ServerName)
} else {
    // koneksi HTTP biasa (tanpa TLS)
}
```

### Field penting di `tls.ConnectionState`

| Field              | Tipe                  | Contoh                          | Arti                  |
| ------------------ | --------------------- | ------------------------------- | --------------------- |
| `Version`          | `uint16`              | `772` (TLS 1.3)                 | Versi TLS yang dipake |
| `CipherSuite`      | `uint16`              | `4865` (TLS_AES_128_GCM_SHA256) | Algoritma enkripsi    |
| `ServerName`       | `string`              | `"localhost"`                   | SNI dari client       |
| `PeerCertificates` | `[]*x509.Certificate` | `[cert.pem]`                    | Certificate chain     |

### Cara Konversi TLS Version

TLS version disimpan sebagai `uint16`, bukan string. Go pake constant:

```go
tls.VersionTLS10 = 0x0301
tls.VersionTLS11 = 0x0302
tls.VersionTLS12 = 0x0303
tls.VersionTLS13 = 0x0304
```

Untuk konversi ke string, pake switch seperti di `server.go`:

```go
func getTLSVersionName(version uint16) string {
    switch version {
    case tls.VersionTLS10:
        return "TLS 1.0"
    case tls.VersionTLS11:
        return "TLS 1.1"
    case tls.VersionTLS12:
        return "TLS 1.2"
    case tls.VersionTLS13:
        return "TLS 1.3"
    }
}
```

---

## Output Server Log

Saat client akses via HTTPS + HTTP/2, log di terminal:

```
received request with HTTP version: HTTP/2
received request with TLS version: TLS 1.3
```

Saat client akses via HTTP (gagal, koneksi ditolak server), tidak ada log karena TLS handshake gagal sebelum handler dipanggil.

---

## Cara Kerja `r.TLS` di Go

```
Client connect ke server:3000
        │
        ▼
TLS handshake terjadi (di level net/http)
        │
        ▼
Kalau handshake sukses → r.TLS = &tls.ConnectionState{...}
Kalau HTTP biasa      → r.TLS = nil
        │
        ▼
Handler dipanggil dengan r (yang udah berisi info TLS)
```

`r.TLS` diisi **oleh `net/http`** sebelum handler dipanggil. Kamu tinggal baca — tidak perlu setup manual.

---

## Catatan

- `r.TLS` hanya ada kalau server pake `ListenAndServeTLS` atau `http.Server{TLSConfig: ...}`
- Kalau server HTTP biasa (`ListenAndServe`), `r.TLS` selalu `nil` meskipun client kirim request HTTPS
- HTTP/2 tidak bisa di-deteksi dari `r.Proto` saja — karena HTTP/2 bisa lewat TLS (h2) atau tanpa TLS (h2c, jarang). Cek `r.TLS` juga untuk konfirmasi.
