# `go get` — What Happens Under the Hood

## TL;DR

```
User runs: go get github.com/someone/lib@v1.2.3
```

1. Resolve module path → cari repo
2. Download source code (zip via proxy atau git clone)
3. Verifikasi checksum (go.sum)
4. Extract ke module cache (`$GOPATH/pkg/mod/`)
5. Update `go.mod` + `go.sum`
6. Compile verification

---

## 1. Resolve Module Path

Go punya urutan prioritas untuk mecari module:

```
GOPROXY = https://proxy.golang.org,direct
```

Artinya:

- **Pertama**: coba ke `proxy.golang.org` (mirror resmi Go)
- **Kedua** (`direct`): kalau proxy gagal atau tidak ditemukan, langsung `git clone` dari repo asli

Go kirim request seperti ini ke proxy:

```
GET https://proxy.golang.org/github.com/someone/lib/@v/v1.2.3.info
```

Proxy balikin informasi version (kalau ada).

### Kapan lewat proxy vs direct?

| Skenario                   | Proxy            | Direct          |
| -------------------------- | ---------------- | --------------- |
| Module public              | ✅ dipake duluan | fallback        |
| `GOPROXY=off`              | ❌               | ❌ error        |
| `GOPRIVATE=*.internal.com` | ❌ dilewati      | ✅ langsung git |
| Proxy offline              | ❌ timeout       | ✅ fallback     |

---

## 2. Download Source Code

### Via Proxy (default)

```
GET https://proxy.golang.org/github.com/someone/lib/@v/v1.2.3.zip
```

Proxy return file `.zip` yang berisi source code module. Isinya cuma file Go yang dibutuhkan — **tidak ada** `.git/`, CI config, README, dll.

### Via Direct (`git clone`)

```bash
git clone --depth 1 --branch v1.2.3 https://github.com/someone/lib.git
```

`--depth 1` = ambil commit terakhir aja, biar cepet. Ini disebut **shallow clone**.

### Hasil Download

File zip disimpan sementara di:

```
$GOPATH/pkg/mod/cache/download/github.com/someone/lib/@v/v1.2.3.zip
```

Beserta file pendamping:

- `v1.2.3.info` — metadata version (waktu, versi)
- `v1.2.3.mod` — `go.mod` dari module itu
- `v1.2.3.ziphash` — hash untuk verifikasi

---

## 3. Checksum Verification (go.sum)

Setelah download, Go verifikasi file:

1. Hitung hash SHA-256 dari file `.zip`
2. Bandingkan dengan hash yang tercatat di `go.sum`
3. Kalau `go.sum` belum ada entry → Go cek ke `sum.golang.org` (checksum database)
4. Hash harus cocok — kalau beda, Go tolak module-nya

```
go.sum entry:
github.com/someone/lib v1.2.3 h1:abc123def456...
```

**Kenapa ini penting?** Mencegah **supply chain attack** — orang jahat mengganti isi module di repo tanpa ketahuan.

### Skipping Checksum Database

Ada kalanya checksum database tidak bisa diakses (misal module private):

```bash
GONOSUMCHECK=github.com/private/* go get ...
```

Atau set `GONOSUMDB` untuk skip database sama sekali.

---

## 4. Extract ke Module Cache

Zip diekstrak ke folder permanen:

```
$GOPATH/pkg/mod/github.com/someone/lib@v1.2.3/
```

Strukturnya:

```
$GOPATH/pkg/mod/
├── cache/                     # cache download (zip, info, mod)
│   └── download/...
├── github.com/
│   └── someone/
│       └── lib@v1.2.3/       # extracted source code
│           ├── lib.go
│           ├── go.mod
│           └── ...
└── sumdb/                     # cache checksum database
```

**Module cache bersifat read-only.** Go bikin file di sini dengan permission `0444` (read-only) untuk mencegah file berubah tanpa sengaja.

### Kenapa perlu cache?

- Download sekali, compile ribuan kali
- Semua project di laptop yang sama pake code yang sama
- Bisa jalan offline (kalau cache sudah ada)

---

## 5. Update `go.mod` dan `go.sum`

### go.mod — Sebelum

```
module my-project

go 1.25.0
```

### go.mod — Sesudah

```
module my-project

go 1.25.0

require github.com/someone/lib v1.2.3
```

### go.sum — Sebelum

(kosong)

### go.sum — Sesudah

```
github.com/someone/lib v1.2.3 h1:abc123...
github.com/someone/lib v1.2.3/go.mod h1:def456...
```

File `go.sum` berisi dua baris per module:

- **Baris 1**: hash dari zip source code
- **Baris 2**: hash dari file `go.mod` module itu sendiri

Keduanya perlu diverifikasi secara terpisah.

---

## 6. Compile Verification

`go get` bukan cuma download — dia juga **compile module yang baru masuk**.

Go jalankan:

1. Parse semua file `.go` di module
2. Resolve semua import (dalam module & ke module lain)
3. Type checking & semantic analysis
4. Pastikan module benar-benar bisa dipake

Kalau gagal compile, Go akan:

```
go: finding module ...
go: downloading ...
go: compiling ...
go: errors found:
# github.com/someone/lib
... compilation error ...

# GO TIDAK ROLLBACK go.mod
# Module tetap tercatat meskipun gagal compile
```

**Penting**: `go get` yang gagal compile **tetap** update `go.mod`. Kamu harus manual hapus baris `require`-nya atau pake `go get ...@none`.

---

## 7. Transitive Dependencies — Kenapa Library Lain Ikut Ke-Install

### Masalah

Kamu `go get A`, eh tiba-tiba `B`, `C`, `D` ikut ke-download. Ini bukan error — ini **transitive dependencies**.

### Analogi

Kamu pesan nasi goreng (A):

- Nasi goreng butuh **nasi** (B)
- Nasi butuh **beras** (C)
- Beras butuh **padi** (D)

Kamu nggak pesan padi, tapi tetap dapat padi. Itu transitive.

### Cara Kerja di Go

```
go.mod A:
require github.com/someone/lib v1.2.3
```

```
go.mod lib (milik github.com/someone/lib):
require (
    golang.org/x/sync v0.1.0
    github.com/autre/autre v2.0.0
)
```

```
go.mod sync:
require (
    golang.org/x/sys v0.5.0
)
```

Maka satu `go get github.com/someone/lib` akan download:

| Module                   | Kenapa?                |
| ------------------------ | ---------------------- |
| `github.com/someone/lib` | langsung di-import     |
| `golang.org/x/sync`      | dependency dari lib    |
| `github.com/autre/autre` | dependency dari lib    |
| `golang.org/x/sys`       | dependency dari x/sync |

Semua masuk ke `go.sum`, tapi **hanya module langsung** yang masuk `require` di `go.mod` kamu:

```
go.mod kamu (setelah go get):
require github.com/someone/lib v1.2.3
```

Sisanya (transitive) dicatat di `go.sum` dan oleh Go akan di-resolve secara otomatis saat compile. Kamu tidak perlu tahu atau peduli — Go urus sendiri.

### Minimum Version Selection (MVS)

Go punya aturan unik: kalau dua module butuh versi berbeda dari dependency yang sama, Go pilih **versi minimum yang mencukupi semua kebutuhan**, bukan versi terbaru.

Contoh:

```
Modul A butuh C v1.0
Modul B butuh C v1.5
```

Go pilih `C v1.5` — karena v1.5 lebih tinggi dari v1.0 dan memenuhi kedua kebutuhan.

Ini beda dengan npm/yarn yang pilih versi terbaru. Go sengaja konservatif biar build-nya stabil dan predictable.

### Kenapa Tidak Semua Library Ikut?

`go get` **hanya download yang diperlukan**. Syaratnya:

- Ada `import` statement yang merujuk ke module itu (langsung atau transitive)
- Atau module itu ada di `go.mod` dan masih diperlukan

Kalau module `A` punya dependency `B` di `go.mod`-nya, tapi nggak ada file `.go` yang `import` `B` → `B` **tidak** ikut di-download. Hanya dependency yang benar-benar digunakan yang di-resolve.

Ini beda dengan `go mod tidy` yang kadang menambahkan dependency baru kalau ada import yang belum tercatat.

### Inspect Dependency Graph

```bash
# lihat semua dependency (termasuk transitive)
go list -m all

# cek kenapa suatu module ada
go mod why github.com/someone/lib

# graph visual
go mod graph | head -20
```

`go mod graph` ngasih output:

```
my-project github.com/someone/lib
github.com/someone/lib golang.org/x/sync
golang.org/x/sync golang.org/x/sys
```

Dari graph ini kamu bisa lihat rantai kenapa `golang.org/x/sys` ikut terinstall — karena `my-project` → `lib` → `x/sync` → `x/sys`.

---

## Diagram Alur

```
User runs: go get github.com/someone/lib@v1.2.3
           │
           ▼
    ┌──────────────┐
    │ 1. Resolve   │─── GOPROXY? GOPRIVATE?
    │    module    │
    └──────┬───────┘
           ▼
    ┌──────────────┐
    │ 2. Download  │─── proxy.zip or git clone
    └──────┬───────┘
           ▼
    ┌──────────────┐
    │ 3. Verify    │─── go.sum + sum.golang.org
    │    checksum  │
    └──────┬───────┘
           ▼
    ┌──────────────┐
    │ 4. Extract   │─── $GOPATH/pkg/mod/
    │    to cache  │
    └──────┬───────┘
           ▼
    ┌──────────────┐
    │ 5. Update    │─── go.mod + go.sum
    │    manifests │
    └──────┬───────┘
           ▼
    ┌──────────────┐
    │ 6. Compile   │─── verification
    │    check     │
    └──────┬───────┘
           ▼
        Done ✅
```

---

## Catatan Tambahan

### GOPROXY Modes

| Value                             | Behavior                             |
| --------------------------------- | ------------------------------------ |
| `https://proxy.golang.org,direct` | default: proxy dulu, fallback direct |
| `direct`                          | langsung git, skip proxy             |
| `off`                             | tidak boleh download module baru     |
| `file:///path/to/local/cache`     | pakai folder lokal sebagai proxy     |

### GONOSUMCHECK vs GONOSUMDB vs GOPRIVATE

| Variable       | Fungsi                                                      |
| -------------- | ----------------------------------------------------------- |
| `GOPRIVATE`    | Module ini: skip proxy, skip sumdb                          |
| `GONOSUMCHECK` | Skip verifikasi checksum (tapi tetap download dari proxy)   |
| `GONOSUMDB`    | Jangan tanya checksum database (tapi bisa verifikasi lokal) |

Biasanya `GOPRIVATE` mencakup kedua fungsi di atas.

### `go get` vs `go mod tidy`

| Perintah             | Fungsi                                                         |
| -------------------- | -------------------------------------------------------------- |
| `go get pkg@version` | Tambah/update dependency spesifik                              |
| `go mod tidy`        | Rapihin `go.mod` — hapus yang tidak dipake, tambah yang kurang |

`go mod tidy` secara implisit juga panggil mekanisme download & verifikasi yang sama kayak `go get`.

### Module Cache di CI/CD

Di lingkungan CI, biasanya set:

```bash
export GOFLAGS=-mod=mod
export GOPROXY=https://proxy.golang.org,direct
```

Atau pakai **Go module mirror** internal perusahaan untuk menghindari rate limiting dari proxy publik.
