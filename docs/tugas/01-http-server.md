# HTTP Server

## Ringkasan
- `package main` + `func main()` = entry point program Go
- `http.HandleFunc(path, handler)` daftarkan handler untuk suatu endpoint
- `http.ListenAndServe(port, nil)` hidupkan server (blocking)
- Port = "pintu" yang bedain satu service dari service lain
- Format `":3000"` = dengerin dari semua network interface

## Tugas

### 1. Cek pemahaman
Jawab tanpa liat kode:
- Kenapa `package main` itu spesial?
- Apa beda `fmt.Println` dan `fmt.Fprintf`?
- Kenapa `err := http.ListenAndServe(...)` harus di-cek?
- Apa yang terjadi kalau port 3000 udah dipake program lain?

### 2. Modifikasi server.go
Tambah satu endpoint baru `/health` yang return `"ok"`, lalu coba akses:
```bash
curl http://localhost:3000/health
```

### 3. Debug
Ganti port ke 80, lalu jalanin server. Error apa yang muncul? Kenapa?

## Tips
- `http.HandleFunc` pake `DefaultServeMux` — cocok buat project kecil
- Error `address already in use` = port dipake program lain
- Port < 1024 butuh `sudo` di Linux
- Selalu handle error dari `ListenAndServe`

## Kalau Masih Bingung
- Baca ulang: `docs/explanations/port.md`
- Baca ulang: `server.go` (komentar tiap baris)

---

## Jawaban

### 1. Cek pemahaman

**Kenapa `package main` itu spesial?**
Karena cuma `package main` yang bisa punya `func main()` — entry point program. Package lain (misal `package handler`) jadi library, bukan executable.

**Apa beda `fmt.Println` dan `fmt.Fprintf`?**
- `fmt.Println` → cetak ke terminal (STDOUT)
- `fmt.Fprintf(w, ...)` → tulis ke writer (response HTTP, file, dll)
- `F` di `Fprintf` = File (writer interface)

**Kenapa `err := http.ListenAndServe(...)` harus di-cek?**
Karena `ListenAndServe` return error kalau gagal (port kepake, permission denied, dll). Kalau tidak di-cek, server mati tanpa notifikasi.

**Apa yang terjadi kalau port 3000 udah dipake program lain?**
Error: `listen tcp :3000: bind: address already in use`. Server gagal jalan.

### 2. Modifikasi server.go
Tambahkan di `func main()`:
```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ok")
})
```
Lalu: `curl http://localhost:3000/health` → `ok`

### 3. Debug port 80
Error: `listen tcp :80: bind: permission denied`. Karena port < 1024 butuh akses root. Solusi: `sudo go run server.go` atau ganti port > 1024.
