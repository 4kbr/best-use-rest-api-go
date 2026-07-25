# HTTP/2, HTTPS, HTTP Connection & TLS Handshake

## HTTP Connection Lifecycle

```
curl -k https://localhost:3000/users
```

1. **DNS Resolve** — localhost → 127.0.0.1 (IPv4) / ::1 (IPv6)
2. **TCP 3-Way Handshake** — SYN → SYN-ACK → ACK (buka koneksi)
3. **TLS Handshake** — negosiasi enkripsi antara curl dan server
4. **ALPN** — negosiasi protokol HTTP (h2 vs http/1.1)
5. **HTTP Request** — curl kirim `GET /users HTTP/2`
6. **HTTP Response** — server balas `200 OK + handling incoming users`

### Diagram

```
curl                          server
  │                              │
  ├── TCP SYN ──────────────────►│
  │◄── TCP SYN-ACK ─────────────┤  TCP
  ├── TCP ACK ──────────────────►│  3-way
  │                              │
  ├── Client Hello ─────────────►│
  │◄── Server Hello ────────────┤
  │◄── Certificate ────────────┤  TLS
  │◄── Server Hello Done ──────┤  Handshake
  ├── Client Key Exchange ──────►│
  ├── Change Cipher Spec ───────►│
  ├── Finished ────────────────►│
  │◄── Change Cipher Spec ─────┤
  │◄── Finished ──────────────┤
  │                              │
  ├── ALPN: h2, http/1.1 ──────►│  ALPN
  │◄── ALPN: h2 ───────────────┤
  │                              │
  ├── GET /users HTTP/2 ────────►│  HTTP
  │◄── 200 OK ─────────────────┤
```

---

## TLS Handshake (4 Langkah)

### 1. Client Hello

```
curl → server: TLS version, cipher suite, random number, ALPN list [h2, http/1.1]
```

Client bilang: "Halo, saya mau koneksi aman. Ini kemampuan saya."

### 2. Server Hello + Certificate

```
server → curl: pilih TLS 1.3, pilih cipher, kirim cert.pem
```

Server bilang: "Ok, kita pake TLS 1.3 dengan cipher ini. Ini identitas saya."

### 3. Key Exchange

```
curl → server: generate pre-master secret, encrypt dengan public key (dari cert), kirim
server: decrypt dengan private key (key.pem)
```

Kedua sisi sekarang punya **session key** yang sama.

### 4. Finished

```
curl  → server: "Change Cipher Spec + Finished" (terenkripsi)
server → curl:  "Change Cipher Spec + Finished" (terenkripsi)
```

Semua data selanjutnya terenkripsi.

### Output curl -v (ringkasan)

```
* TLSv1.3 (OUT), Client hello
* TLSv1.3 (IN),  Server hello
* TLSv1.3 (IN),  Certificate
* TLSv1.3 (IN),  Finished
* TLSv1.3 (OUT), Change cipher + Finished
* SSL connection using TLSv1.3 / TLS_AES_128_GCM_SHA256
```

---

## TLS vs SSL

| Aspek          | SSL                         | TLS                               |
| -------------- | --------------------------- | --------------------------------- |
| Kepanjangan    | Secure Sockets Layer        | Transport Layer Security          |
| Versi terakhir | SSL 3.0 (1996)              | TLS 1.3 (2018)                    |
| Status         | ❌ Deprecated (semua versi) | ✅ Standard saat ini              |
| Keamanan       | Rawan POODLE, BEAST attack  | Aman (kecuali 1.0/1.1 deprecated) |

Istilah "SSL" masih sering dipakai, tapi yang dimaksud sebenarnya TLS.
Contoh: `ListenAndServeTLS` — namanya TLS, bukan SSL.

---

## h2 vs h2c

HTTP/2 punya dua mode:

| Mode    | Kepanjangan      | TLS      | Browser          | gRPC                 |
| ------- | ---------------- | -------- | ---------------- | -------------------- |
| **h2**  | HTTP/2 over TLS  | ✅ Wajib | ✅ Support       | ✅ Default           |
| **h2c** | HTTP/2 cleartext | ❌ Tidak | ❌ Tidak support | ⚠️ Hanya development |

Server kita pake **h2** — HTTP/2 di atas TLS. Inilah yang dipakai browser dan gRPC.

---

## ALPN — Protocol Negotiation

ALPN terjadi di **dalam** TLS handshake (setelah Certificate, sebelum Finished).

```
Client: "Saya bisa h2, http/1.1"
Server: "Saya pilih h2"
```

| Skenario    | Round Trip                                  |
| ----------- | ------------------------------------------- |
| Dengan ALPN | 0 tambahan — langsung ke HTTP/2             |
| Tanpa ALPN  | +1 round trip — HTTP/1.1 dulu, lalu upgrade |

---

## HTTP/1.1 vs HTTP/2

### HTTP/1.1

```
Koneksi 1: [ GET /orders  ]──►[ response ]
Koneksi 2: [ GET /users   ]──►[ response ]
Koneksi 3: [ GET /products]──►[ response ]
```

- 1 request = 1 koneksi (sequential)
- Head-of-line blocking: request berikutnya nunggu response sebelumnya
- Header plaintext tiap request (boros bandwidth)

### HTTP/2

```
1 koneksi: ┌─ GET /orders ──► response ─┐
           ├─ GET /users  ──► response ─┤
           └─ GET /products► response ─┘
```

- **Multiplexing**: banyak request dalam 1 koneksi
- **Binary protocol**: parsing lebih ringan
- **HPACK**: header dikompres, cuma dikirim sekali

### Tabel

| Aspek                 | HTTP/1.1               | HTTP/2               |
| --------------------- | ---------------------- | -------------------- |
| Format                | Teks (ASCII)           | Binary               |
| Request per koneksi   | 1 (sequential)         | Banyak (multiplexed) |
| Header efficiency     | Plaintext tiap request | Compressed (HPACK)   |
| Head-of-line blocking | Ada                    | Tidak ada            |
| TLS requirement       | Opsional               | Wajib (di browser)   |

---

## gRPC & HTTP/2

gRPC menggunakan HTTP/2 sebagai transport layer-nya.

### Kenapa gRPC pilih HTTP/2?

- **Multiplexing** — banyak RPC call dalam 1 koneksi TCP
- **Binary framing** — efisien untuk protobuf (binary, bukan JSON)
- **Server push** — server bisa kirim data tanpa diminta client
- **Streaming** — client & server bisa kirim data kapan saja (bidirectional)

### Hubungan dengan project ini

Server kita sudah enable HTTP/2 via `http2.ConfigureServer()`. Ini persiapan untuk gRPC — karena gRPC butuh HTTP/2 sebagai fondasi.

---

## Go Implementation

```go
tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}

server := &http.Server{
    Addr:      ":3000",
    TLSConfig: tlsConfig,
}

http2.ConfigureServer(server, &http2.Server{})

server.ListenAndServeTLS("cert.pem", "key.pem")
```

### Di Handler

```go
r.Proto              // "HTTP/1.1" atau "HTTP/2"
r.TLS                // nil (HTTP) / tidak nil (HTTPS)
r.TLS.Version        // 0x0304 = TLS 1.3
r.TLS.CipherSuite    // 4865 = TLS_AES_128_GCM_SHA256
```

---

## Testing dengan Curl

| Command                          | Hasil                                    |
| -------------------------------- | ---------------------------------------- |
| `curl http://...`                | ❌ 400 Bad Request (HTTPS server)        |
| `curl https://...`               | ❌ SSL certificate problem (self-signed) |
| `curl -k https://...`            | ✅ Berhasil (skip verify)                |
| `curl --http2 -k -v https://...` | ✅ Debug lengkap                         |

---

## Best Practice

- `cert.pem` → publik, boleh di-commit
- `key.pem` → **rahasia**, jangan pernah di-commit
- Self-signed cert → development saja
- Production → Let's Encrypt (gratis) atau CA berbayar
- TLS 1.0/1.1 deprecated → selalu `MinVersion: tls.VersionTLS12`
- HTTP/2 otomatis via ALPN → tidak perlu setup tambahan
- TLS 1.3 lebih cepat dari 1.2 (1 round trip vs 2)
