# Go Module (go get & go mod tidy)

## Ringkasan
- `go get pkg@version` → download dependency + update `go.mod`
- `go get` otomatis download **transitive dependencies** (dependency dari dependency)
- `go mod tidy` → bersihin `go.mod` (hapus yg tak dipake, tambah yg kurang)
- `go.sum` berisi hash untuk verifikasi keamanan module
- `GOPROXY`, `GOPRIVATE`, `GONOSUMCHECK` kontrol darimana module di-download

## Tugas

### 1. Cek pemahaman
- Apa beda `go get` dan `go mod tidy`?
- Kenapa kadang `go get A` ikut download B, C, D?
- Apa fungsi `go.sum`? Kenapa perlu di-commit?
- Apa yang dimaksud MVS (Minimum Version Selection)?

### 2. Praktik
Cek dependency graph project ini:
```bash
go mod graph
```
Identifikasi: module apa saja yang terinstall dan siapa yang membutuhkan siapa.

### 3. Skenario
Kamu hapus baris `import "golang.org/x/net/http2"` dari `server.go`, lalu jalankan `go mod tidy`. Apa yang terjadi pada `go.mod`?

### 4. Tebak
Kalau `GOPROXY=off` terus kamu jalanin `go get github.com/gin-gonic/gin`, kira-kira error apa yang muncul?

## Tips
- Jangan edit `go.mod` manual — biarin Go yang urus
- Jalankan `go mod tidy` sebelum commit
- `go.sum` jangan dihapus — isinya penting buat verifikasi
- Kalau ada error `missing go.sum entry` → tinggal `go mod tidy`

## Kalau Masih Bingung
- Baca ulang: `docs/explanations/go-get.md`
- Baca ulang: `docs/explanations/go-mod-tidy.md`

---

## Jawaban

### 1. Cek pemahaman

**Apa beda `go get` dan `go mod tidy`?**
- `go get` → tambah/update dependency spesifik. Cocok pas butuh library baru.
- `go mod tidy` → sinkronisasi semua: hapus yang tak dipakai + tambah yang kurang.

**Kenapa kadang `go get A` ikut download B, C, D?**
Karena **transitive dependencies**. A butuh B, B butuh C, C butuh D. Go download semua yang diperlukan biar A bisa dipake.

**Apa fungsi `go.sum`? Kenapa perlu di-commit?**
`go.sum` berisi hash SHA-256 dari setiap module. Fungsinya **verifikasi** — memastikan module yang di-download hari ini sama persis dengan kemarin. Mencegah supply chain attack. Wajib di-commit bareng `go.mod`.

**Apa yang dimaksud MVS (Minimum Version Selection)?**
Aturan Go: kalau dua dependency butuh versi berbeda dari module yang sama, Go pilih **versi minimum yang cukup** — bukan yang terbaru. Contoh: A butuh C v1.0, B butuh C v1.5 → Go pilih v1.5.

### 2. Praktik

```
go mod graph
simpleapi golang.org/x/net@v0.57.0
golang.org/x/net@v0.57.0 golang.org/x/text@v0.40.0
```

Artinya: `simpleapi` → butuh `golang.org/x/net` → butuh `golang.org/x/text`.

### 3. Skenario
`go mod tidy` akan hapus baris `require golang.org/x/net` dari `go.mod` karena tidak ada lagi file yang import package `http2`. `golang.org/x/text` juga ikut hilang karena transitive dan tidak diperlukan lagi.

### 4. Tebak
Error: `go: module github.com/gin-gonic/gin: Get "https://proxy.golang.org/...": GOPROXY list is not allowed` atau `GOPROXY=off is not allowed`. Karena `GOPROXY=off` = jangan download apapun dari mana pun.
