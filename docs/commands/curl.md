# Curl — Testing HTTPS & HTTP/2

Dokumen ini mencatat hasil curl saat ngakses server HTTPS, dari gagal sampai berhasil.

---

## Hasil 1: Gagal — HTTP ke HTTPS Server

```bash
curl -v http://localhost:3000/users
```

Output:

```
* Host localhost:3000 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:3000...
* Connected to localhost (::1) port 3000
> GET /users HTTP/1.1
> Host: localhost:3000
> User-Agent: curl/8.5.0
> Accept: */*
>
* HTTP 1.0, assume close after body
< HTTP/1.0 400 Bad Request
<
Client sent an HTTP request to an HTTPS server.
* Closing connection
```

### Kenapa gagal?

| Penyebab                 | Penjelasan                                        |
| ------------------------ | ------------------------------------------------- |
| URL pake `http://`       | Curl kirim request **HTTP biasa** (plaintext)     |
| Server pake HTTPS        | Server hanya terima koneksi **TLS** (terenkripsi) |
| Server baca byte pertama | Yang dapet: `G` dari `GET` — bukan TLS handshake  |

Server jawab: **`400 Bad Request`** — "Client sent an HTTP request to an HTTPS server."

**Fix:** Ganti `http://` jadi `https://`.

---

## Hasil 2: Gagal — HTTPS tanpa `-k`

```bash
curl -v https://localhost:3000/users
```

Output:

```
* Host localhost:3000 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:3000...
* Connected to localhost (::1) port 3000
* ALPN: curl offers h2,http/1.1
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
*  CAfile: /etc/ssl/certs/ca-certificates.crt
*  CApath: /etc/ssl/certs
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (OUT), TLS alert, unknown CA (560):
* SSL certificate problem: self-signed certificate
* Closing connection
curl: (60) SSL certificate problem: self-signed certificate
```

### Kenapa gagal?

| Step                          | Apa yang terjadi                                     |
| ----------------------------- | ---------------------------------------------------- |
| `TLS handshake, Client hello` | Curl bilang "mau koneksi TLS" ke server              |
| `TLS handshake, Server hello` | Server jawab "ok, ini TLS 1.3"                       |
| `TLS handshake, Certificate`  | Server kirim cert.pem ke curl                        |
| `TLS alert, unknown CA`       | Curl cek cert → **self-signed**, bukan dari CA resmi |
| `SSL certificate problem`     | Curl putuskan koneksi — cert tidak dikenal           |

Curl punya daftar CA terpercaya (di `/etc/ssl/certs/`). Cert self-signed tidak ada di daftar itu, jadi curl curiga dan menolak.

**Fix:** Flag `-k` (atau `--insecure`) — skip verifikasi cert. Aman untuk development.

---

## Hasil 3: Berhasil — HTTPS + `-k`

```bash
curl -v -k https://localhost:3000/users
```

Output:

```
* Host localhost:3000 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:3000...
* Connected to localhost (::1) port 3000
* ALPN: curl offers h2,http/1.1
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS_AES_128_GCM_SHA256 / X25519 / RSASSA-PSS
* ALPN: server accepted h2
* Server certificate:
*  subject: C=AU; ST=Some-State; O=Internet Widgits Pty Ltd
*  start date: Jul 25 07:22:25 2026 GMT
*  expire date: Jul 25 07:22:25 2027 GMT
*  issuer: C=AU; ST=Some-State; O=Internet Widgits Pty Ltd
*  SSL certificate verify result: self-signed certificate (18), continuing anyway.
*   Certificate level 0: Public key type RSA (2048/112 Bits/secBits), signed using sha256WithRSAEncryption
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* using HTTP/2
* [HTTP/2] [1] OPENED stream for https://localhost:3000/users
* [HTTP/2] [1] [:method: GET]
* [HTTP/2] [1] [:scheme: https]
* [HTTP/2] [1] [:authority: localhost:3000]
* [HTTP/2] [1] [:path: /users]
* [HTTP/2] [1] [user-agent: curl/8.5.0]
* [HTTP/2] [1] [accept: */*]
> GET /users HTTP/2
> Host: localhost:3000
> User-Agent: curl/8.5.0
> Accept: */*
>
< HTTP/2 200
< content-type: text/plain; charset=utf-8
< content-length: 23
< date: Sat, 25 Jul 2026 14:55:24 GMT
<
* Connection #0 to host localhost left intact
handling incoming users
```

### Kenapa berhasil?

| Baris penting                                     | Arti                                                       |
| ------------------------------------------------- | ---------------------------------------------------------- |
| `ALPN: server accepted h2`                        | Server setuju pake HTTP/2                                  |
| `SSL connection using TLSv1.3 / ...`              | Koneksi aman: **TLS 1.3** + cipher modern                  |
| `self-signed certificate (18), continuing anyway` | Curl tau cert self-signed, tapi karena `-k`, dia lanjutkan |
| `using HTTP/2`                                    | Koneksi pake **HTTP/2**, bukan HTTP/1.1                    |
| `GET /users HTTP/2`                               | Request dikirim pake protokol HTTP/2                       |
| `HTTP/2 200`                                      | Response 200 via HTTP/2                                    |

### Cara baca output curl verbose

Output curl `-v` terbagi jadi 3 bagian:

```
* baris dengan *    → informasi koneksi (TLS, DNS, cert)
> baris dengan >    → header request (dari client ke server)
< baris dengan <    → header response (dari server ke client)
```

---

## Ringkasan

| Skenario            | Command                                  | Hasil                        |
| ------------------- | ---------------------------------------- | ---------------------------- |
| HTTP ke HTTPS       | `curl http://localhost:3000`             | ❌ `400 Bad Request`         |
| HTTPS tanpa trust   | `curl https://localhost:3000`            | ❌ `SSL certificate problem` |
| HTTPS + skip verify | `curl -k https://localhost:3000`         | ✅ Berhasil                  |
| HTTPS + HTTP/2      | `curl --http2 -k https://localhost:3000` | ✅ HTTP/2                    |

**Catatan:** Di kode kita, HTTP/2 sudah aktif via `http2.ConfigureServer()`. Curl otomatis negosiasi HTTP/2 selama koneksi HTTPS — tidak perlu flag `--http2` (tapi bisa dikasih buat mastiin).
