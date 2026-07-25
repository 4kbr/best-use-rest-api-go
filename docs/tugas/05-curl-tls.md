# Curl, TLS & HTTP Request Inspection

## Ringkasan
- `curl http://` → gagal karena server HTTPS (400 Bad Request)
- `curl https://` → gagal karena cert self-signed (SSL certificate problem)
- `curl -k https://` → berhasil (`-k` skip verifikasi cert)
- `r.Proto` = HTTP version dari request
- `r.TLS` = nil (HTTP) / tidak nil (HTTPS), berisi info TLS

## Tugas

### 1. Cek pemahaman
Jawab tanpa liat kode:
- Apa isi `r.Proto` kalau client pake HTTP/2?
- Apa beda output curl baris `>` dan `<`?
- Kenapa `r.TLS` bisa `nil` meskipun server kita HTTPS?
- Di output curl yang berhasil, cari baris yang nunjukkin koneksi pake HTTP/2. Baris apa itu?

### 2. Debug
Jalankan server, lalu akses dengan cara yang salah:
```bash
# Tanpa flag -k
curl https://localhost:3000/orders
```
Apa error code curl yang keluar? (angka dalam kurung, misal `(60)`)

### 3. Praktik
Jalankan server, lalu akses dengan:
```bash
curl -k https://localhost:3000/orders -v
```
Perhatikan log di terminal server. Catat:
- HTTP version apa yang tercetak?
- TLS version apa yang tercetak?

### 4. Modifikasi
Comment dulu baris `http2.ConfigureServer()` di `server.go`, restart server, lalu akses lagi dengan:
```bash
curl -k https://localhost:3000/orders -v
```

Apa yang berubah di:
- output curl? (cari baris `using HTTP/...`)
- log di terminal server? (cari baris `received request with HTTP version:`)

### 5. Eksplorasi
Coba akses dengan protokol HTTP/1.1 secara paksa:
```bash
curl --http1.1 -k https://localhost:3000/orders -v
```
Lalu bandingkan output dengan yang HTTP/2. Apa perbedaan yang terlihat?

## Tips
- `-v` (verbose) di curl = mode debug — tunjukin semua detail koneksi
- `-k` = insecure — skip verifikasi certificate. **Hanya untuk development**
- Kalau lupa `-k`, error `(60) SSL certificate problem`
- Angka error curl seperti `(60)` bisa di-search di web untuk detail
- `r.Proto` dan `r.TLS` cuma bisa dibaca di dalem handler — request sudah lolos TLS handshake

## Kalau Masih Bingung
- Baca ulang: `docs/commands/curl.md`
- Baca ulang: `docs/explanations/http-request-inspection.md`
- Baca ulang: `server.go` (fungsi `logRequestDetails`)

---

## Jawaban

### 1. Cek pemahaman

**Apa isi `r.Proto` kalau client pake HTTP/2?**
`"HTTP/2"` — bukan `"HTTP/2.0"`.

**Apa beda output curl baris `>` dan `<`?**
- `>` = header request (dari client ke server)
- `<` = header response (dari server ke client)
- `*` = informasi koneksi (TLS, DNS, cert)

**Kenapa `r.TLS` bisa `nil` meskipun server kita HTTPS?**
`r.TLS` diisi per request. Kalau client connect tapi pakai HTTP biasa (bukan HTTPS), TLS handshake tidak terjadi. Tapi di kasus kita, request HTTP ditolak oleh server sebelum handler — jadi `r.TLS` tidak relevan. Di koneksi HTTPS yang sukses, `r.TLS` pasti tidak nil.

**Di output curl yang berhasil, cari baris yang nunjukkin koneksi pake HTTP/2. Baris apa itu?**
```
* using HTTP/2
```
Ada juga:
```
* ALPN: server accepted h2
```
Yang berarti server setuju pake HTTP/2.

### 2. Debug
Error code: `(60)` — `SSL certificate problem: self-signed certificate`.

### 3. Praktik
```
received request with HTTP version: HTTP/2
received request with TLS version: TLS 1.3
```

### 4. Modifikasi (tanpa http2.ConfigureServer)

**Output curl:**
```
* using HTTP/1.1
```
Atau baris ALPN: `* ALPN: server did not accept h2` atau sama sekali tidak negosiasi HTTP/2.

**Log server:**
```
received request with HTTP version: HTTP/1.1
received request with TLS version: TLS 1.3
```
HTTP version turun dari `HTTP/2` ke `HTTP/1.1`.

### 5. Eksplorasi
- Dengan `--http1.1`: koneksi pake HTTP/1.1 walau server support HTTP/2
- Server log: `received request with HTTP version: HTTP/1.1`
- Output curl: `* using HTTP/1.1`
- Tidak ada baris `* [HTTP/2] [...]` — karena koneksi bukan HTTP/2
- Response tetap 200, cuma protokolnya berbeda
