# HTTP/2, HTTPS & TLS

## Ringkasan
- HTTPS = HTTP + TLS (enkripsi data antara client-server)
- TLS pakai **certificate** (public) + **private key** (rahasia)
- Self-signed cert: buat sendiri pakai openssl, cukup untuk development
- **HTTP/2 butuh TLS** — browser cuma support h2, bukan h2c
- `http.ListenAndServe` → HTTP/1.1 (tanpa enkripsi)
- `server.ListenAndServeTLS(cert, key)` → HTTPS (dengan enkripsi)
- `http2.ConfigureServer(server, nil)` → enable HTTP/2

## Tugas

### 1. Cek pemahaman
Jawab tanpa liat kode:
- Kenapa `key.pem` tidak boleh di-commit ke git?
- Apa fungsi `MinVersion: tls.VersionTLS12`?
- Kenapa di kode kita pake `&http.Server{}` bukan `http.ListenAndServe`?
- Apa beda `ListenAndServe` dan `ListenAndServeTLS`?
- Kenapa curl perlu flag `-k` saat akses server kita?

### 2. Debug
Jalankan server, lalu coba akses dengan:
```bash
# Format salah — HTTP ke server HTTPS
curl http://localhost:3000/orders
```
Error apa yang muncul? Catat dan jelaskan.

### 3. Verifikasi HTTP/2
Jalankan:
```bash
curl --http2 -k https://localhost:3000/orders -v
```
Cari baris `* Using HTTP/2`. Kalau tidak muncul, apa yang salah?

### 4. Skenario
Hapus baris `http2.ConfigureServer()` dari kode, restart server, lalu akses dengan `curl --http2 -k`. Apa yang terjadi?

## Tips
- Self-signed cert = ok untuk development, **jangan** dipakai di production
- Untuk production: pakai Let's Encrypt (free) atau CA berbayar
- TLS 1.0 dan 1.1 sudah deprecated — selalu set `MinVersion: tls.VersionTLS12`
- Kalau lupa generate cert: `openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 365`

## Kalau Masih Bingung
- Baca ulang: `docs/commands/certificate.md`
- Baca ulang: `server.go` (bagian TLS, cert, HTTP/2)

---

## Jawaban

### 1. Cek pemahaman

**Kenapa `key.pem` tidak boleh di-commit ke git?**
Karena `key.pem` adalah **private key** — kunci rahasia untuk mendekripsi data. Kalau bocor ke publik (apalagi di GitHub), orang lain bisa menyamar sebagai server kita.

**Apa fungsi `MinVersion: tls.VersionTLS12`?**
Mencegah koneksi menggunakan TLS versi lama (1.0, 1.1) yang sudah terbukti punya celah keamanan. Set minimal ke 1.2.

**Kenapa di kode kita pake `&http.Server{}` bukan `http.ListenAndServe`?**
Karena kita perlu ngasih `TLSConfig` dan enable HTTP/2 lewat `http2.ConfigureServer()`. `http.ListenAndServe` tidak bisa dikustomisasi — cuma pake `DefaultServeMux` tanpa TLS.

**Apa beda `ListenAndServe` dan `ListenAndServeTLS`?**
- `ListenAndServe` → HTTP biasa (tanpa enkripsi)
- `ListenAndServeTLS(cert, key)` → HTTPS (dengan enkripsi TLS)

**Kenapa curl perlu flag `-k` saat akses server kita?**
Karena certificate kita **self-signed** — tidak ditandatangani oleh CA resmi. Curl secara default nolak koneksi ke server dengan cert tidak dikenal. `-k` = skip verifikasi.

### 2. Debug

```
curl: (35) error:0A000438:SSL routines::tlsv1 alert internal error
```

Atau error serupa. Karena client HTTP ngirim request plaintext, tapi server mengharapkan TLS handshake. Server bingung dan putuskan koneksi.

### 3. Verifikasi HTTP/2
Output harus ada:
```
* Using HTTP/2
```
Kalau tidak muncul, kemungkinan:
- Lupa panggil `http2.ConfigureServer()`
- Curl versi lama tidak support HTTP/2
- Flag `--http2` tidak dikasih

### 4. Skenario
Tanpa `http2.ConfigureServer()`, server cuma HTTP/1.1. Curl dengan `--http2` akan fallback ke HTTP/1.1 — tidak error, tapi kecepatan HTTP/2 hilang. Output curl: `* Using HTTP/1.1`.
