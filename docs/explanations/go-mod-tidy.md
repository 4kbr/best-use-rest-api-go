# `go mod tidy` — What Happens Under the Hood

## TL;DR

`go mod tidy` merapikan `go.mod` dan `go.sum`:

1. **Tambah** module yang di-import di kode tapi belum tercatat
2. **Hapus** module yang tercatat tapi tidak dipakai
3. **Update** `go.sum` agar cocok dengan isi `go.mod`

---

## Kenapa Perlu `go mod tidy`

`go.mod` bisa kotor karena:

| Penyebab                       | Akibat                                   |
| ------------------------------ | ---------------------------------------- |
| `go get` lalu tidak jadi pakai | module orphans                           |
| Ganti library                  | module lama masih tercatat               |
| Hapus file yang `import`       | dependency tidak lagi diperlukan         |
| Merge conflict di `go.mod`     | isi `go.mod` tidak konsisten dengan kode |

`go mod tidy` adalah **vacuum cleaner** untuk `go.mod` dan `go.sum`.

---

## 3 Hal yang Dilakukan `go mod tidy`

### 1. Tambah Missing Dependencies

Go scan semua file `.go` di project (kecuali yang `_test.go` dengan build tags tertentu).

Kalau ada import:

```go
import "github.com/someone/lib"
```

Tapi `go.mod` belum punya `require github.com/someone/lib` → Go tambahkan otomatis.

Yang ke-scan:

- Semua file `*.go` di project root
- Semua file `*.go` di semua sub-package (dalam folder)
- File test (`*_test.go`) **juga di-scan** — dependency untuk test juga ikut

### 2. Hapus Unused Dependencies

Kalau `go.mod` punya:

```
require github.com/old/lib v1.0.0
```

Tapi tidak ada satupun file `.go` yang `import "github.com/old/lib"` → Go hapus baris `require`-nya.

**Penting**: Go tetap hitung transitive dependency. Misal:

```
go.mod kamu require A
A butuh B
B butuh C
Kamu tidak import C langsung
```

→ `C` tetap dianggap **used** karena diperlukan oleh `B` yang diperlukan oleh `A`.

Jadi yang dihapus cuma module yang benar-benar **tidak tersentuh sama sekali** dalam rantai dependency.

### 3. Sinkronisasi `go.sum`

`go.sum` harus cocok dengan `go.mod`. `go mod tidy`:

1. Hitung semua module yang diperlukan (langsung + transitive)
2. Download yang belum ada di cache
3. Generate hash untuk setiap module
4. Hapus entry `go.sum` yang tidak diperlukan
5. Tambah entry `go.sum` yang kurang

Hasilnya `go.sum` cuma berisi module yang relevan — tidak lebih, tidak kurang.

---

## Mode: `-compat`

```bash
go mod tidy -compat=1.25.0
```

Menentukan Go version yang dipakai untuk resolve dependency graph. Berguna kalau project mu harus kompatibel dengan Go version tertentu.

Kalau tidak dikasih flag, `go mod tidy` pakai version yang tercantum di `go.mod` (baris `go 1.xx.x`).

Kapan perlu `-compat`:

- Project punya multiple version Go
- Mau pastikan module yang dipilih bekerja di Go version lama
- Upgrade Go version tapi belum siap commit `go.mod`

---

## `go mod tidy` vs `go get`

| Perintah      | Fungsi                        | ngaruh ke `go.mod`             |
| ------------- | ----------------------------- | ------------------------------ |
| `go get pkg`  | Tambah/update satu dependency | nambah/ubah baris `require`    |
| `go mod tidy` | Sinkronisasi semua dependency | nambah + hapus baris `require` |

Urutan yang recommended:

```bash
# pas develop
go get github.com/someone/lib

# sebelum commit — bersihin
go mod tidy
```

`go mod tidy` bisa **nambah** dependency baru kalau ada import yang belum terdaftar, tapi dia nggak nambah versi spesifik — dia pilih **versi terendah yang compatible** (sesuai MVS).

---

## Kapan Harus `go mod tidy`

| Skenario                                     | Wajib `go mod tidy`?                     |
| -------------------------------------------- | ---------------------------------------- |
| Baru `git pull` dan ada conflict di `go.mod` | ✅ ya                                    |
| Selesai nambah fitur, sebelum commit         | ✅ recommended                           |
| Ganti library (A ke B)                       | ✅ ya, B harus ditambah, A harus dihapus |
| `go build` gagal karena missing module       | ✅ coba dulu                             |
| Setiap kali `go run`                         | ❌ tidak perlu, tapi nggak ada ruginya   |
| Project stabil, tidak ada perubahan          | ❌ skip                                  |

---

## Cara Kerja Detail

```
Step 1: Scan semua file .go (termasuk _test.go)
         ↓
Step 2: Kumpulkan semua import statements
         ↓
Step 3: Bandingkan dengan require di go.mod
         ↓
Step 4: Hapus require yang tidak di-import
         ↓
Step 5: Tambah require yang di-import tapi belum ada
         ↓
Step 6: Download module baru (kalau perlu)
         ↓
Step 7: Generate go.sum yang cocok
         ↓
Step 8: Tulis go.mod + go.sum
```

### Detail Step 1 — Scanning

Yang di-scan:

```
project/
├── main.go              ✅
├── helper.go            ✅
├── helper_test.go       ✅ (test file tetap di-scan)
├── internal/
│   └── handler.go       ✅ (sub-package)
├── cmd/
│   └── server/main.go   ✅
└── vendor/              ❌ (skip, pakai vendor mode)

Build tags:
File dengan `//go:build ignore`  ❌ dilewati
File dengan `//go:build !linux`  ✅ tergantung OS
```

### Detail Step 4-5 — Penambahan & Penghapusan

Sebelum `go mod tidy`:

```
module my-project
go 1.25.0
require (
    github.com/A v1.0.0    ← masih dipakai ✅
    github.com/B v1.0.0    ← sudah tidak di-import ❌
    github.com/C v1.0.0    ← dipakai via import A ❓
)
```

Setelah `go mod tidy`:

```
module my-project
go 1.25.0
require github.com/A v1.0.0    ← B dihapus, C tetap karena A butuh C
```

---

## Interaksi dengan Vendor

Proyek yang pakai `vendor/` perlu step tambahan:

```bash
go mod tidy          # bersihin go.mod + go.sum
go mod vendor        # rebuild vendor/ folder sesuai go.mod baru
```

Kadang developer lupa `go mod vendor` setelah `go mod tidy`, dan CI/CD yang pakai `-mod=vendor` jadi error.

---

## Best Practice

1. **Jalankan `go mod tidy` sebelum commit** — jaga `go.mod` tetap bersih
2. **Commit `go.mod` dan `go.sum` bersama-sama** — keduanya harus sinkron
3. **Jangan edit `go.mod` manual** — Go bisa nulis ulang dan formatting-nya berubah
4. **Kalau conflict `go.mod`**, resolve dulu baru `go mod tidy`

---

## Error yang Sering Muncul

| Error                                                     | Penyebab                       | Solusi                                 |
| --------------------------------------------------------- | ------------------------------ | -------------------------------------- |
| `go: updates to go.mod needed, disabled by -mod=readonly` | `go.mod` read-only karena flag | `go mod tidy` atau `go build -mod=mod` |
| `missing go.sum entry`                                    | `go.sum` ketinggalan isi       | `go mod tidy`                          |
| `module declared in go.mod but not required`              | conflict atau edit manual      | `go mod tidy`                          |

Semua error di atas solusinya satu: **`go mod tidy`**.
