# AGENTS.md — Best Practice API

## Profil Project

| Item | Detail |
|---|---|
| Nama | `best-practice-api` |
| Stack | Go, std lib `net/http` |
| Module | `best-practice-api` |
| Go version | 1.25.0 |


Project ini adalah REST API sederhana sebagai learning ground untuk best practice Go. Semua kode harus idiomatik, mudah dibaca, dan mudah dikembangkan.

---

## Struktur Direktori

Ikuti konvensi ini saat menambah file/folder:

```
.
├── cmd/
│   └── server/          # entry point → func main()
├── internal/            # private package, tidak bisa di-import dari luar
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # middleware (logging, auth, CORS, dll)
│   ├── model/           # struct/type definitions
│   ├── repository/      # data access layer
│   └── service/         # business logic
├── pkg/                 # public package, bisa di-import dari luar
├── go.mod
└── server.go            # (sementara, nanti pindah ke cmd/server/)
```

Aturan:
- `internal/` untuk kode yang **tidak boleh** diakses dari luar module.
- `pkg/` untuk kode yang **sengaja** diekspos.
- Setiap package di `internal/` punya tanggung jawab **satu hal**.

---

## Go Best Practices

### Package Naming
- Lowercase, satu kata.
- Hindari `utils`, `common`, `helper`. Cari nama deskriptif.
- Contoh: `handler`, `repo`, `middleware`, `validator`.

### File Naming
- `snake_case.go` — contoh: `user_handler.go`, `auth_middleware.go`.
- Satu file untuk satu tipe utama, kecuali memang kecil.

### Variable & Function Naming
- **Short-lived variable** → nama pendek: `i`, `r`, `w`, `err`.
- **Exported function/type** → `PascalCase`.
- **Unexported** → `camelCase`.
- **Acronym** → semua kapital: `HTTP`, `URL`, `API`, `ID`.

### Error Handling
- **Always check errors.** Jangan pakai `_` untuk error.
- Return error ke caller; jangan `panic`/`log.Fatal` di dalam function non-main.
- Gunakan `fmt.Errorf("context: %w", err)` untuk wrapping.
- Buat custom error type jika perlu informasi tambahan.

### Idiomatic Go
- Accept interfaces, return structs.
- Gunakan zero-value initialization (`var s Struct`) daripada `new()`.
- Prefer `range` daripada index loop.
- Jangan export field/data tanpa alasan jelas.
- Gunakan `http.Handler` interface — buat handler method dengan `ServeHTTP` atau handler function.

### Imports
- Kelompokkan: std lib → external → internal.
- Pisahkan tiap grup dengan baris kosong.

```go
import (
    "context"
    "fmt"
    "net/http"

    "github.com/gorilla/mux"

    "best-practice-api/internal/model"
)
```

### Testing
- Package: `testing`.
- File: `xxx_test.go` di samping file yang di-test.
- Nama fungsi: `TestXxx` (PascalCase).
- Gunakan **table-driven test** untuk multiple cases.
- Test harus **deterministic** dan **independen**.
- Jangan test external dependency langsung — gunakan interface + mock.
- Jalankan `go test ./...` sebelum commit.

---

## REST API Patterns

### Routing
- Gunakan **method-based routing**:
  - `GET /api/v1/resource` → list/get
  - `POST /api/v1/resource` → create
  - `PUT /api/v1/resource/:id` → update
  - `DELETE /api/v1/resource/:id` → delete
- Prefix `api/v1/` untuk versioning.

### Response Envelope
Semua response JSON pakai format konsisten:

```go
// Success
{
    "data": { ... },
    "meta": {
        "page": 1,
        "total": 42
    }
}

// Error
{
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "field 'email' is required"
    }
}
```

### HTTP Status Codes
- `200` — sukses (GET, PUT)
- `201` — created (POST)
- `204` — no content (DELETE)
- `400` — bad request / validation error
- `401` — unauthorized
- `403` — forbidden
- `404` — not found
- `409` — conflict
- `422` — unprocessable entity
- `500` — internal server error

---

## Agent Rules

### Sebelum Ngoding
1. Baca file yang relevan dulu — pahami kode yang sudah ada.
2. Cari tahu konvensi yang udah dipakai (naming, struktur, pattern).
3. Tanya sebelum:
   - Nambah dependency eksternal.
   - Refactor besar.
   - Hapus kode yang sudah ada.

### Saat Ngoding
4. Ikuti gaya kode yang sudah ada — jangan mixing style.
5. Jangan hardcode config — pakai env variable.
6. Jangan tambah komentar yang nggak perlu — kode harus self-documenting.
7. Setiap operasi I/O harus handle error.
8. Jangan expose secret/key di kode.

### Sebelum Selesai
9. Jalankan `go vet ./...` — cek masalah statis.
10. Jalankan `go test ./...` — pastikan semua test pass.
11. Cek ulang diff — pastikan hanya file yang dimaksud yang berubah.

### Guard for agents
- jangan pernah buat commit sendiri, selalu minta persetujuan user
- jangan baca .env sendiri, selalu minta persetujuan user
- sebelum mengerjakan setiap tugas, selalu baca AGENTS.md sebagai 


---

## Commit Convention

Gunakan **Conventional Commits**:

```
<type>: <deskripsi singkat>

<opsional: body detail>
```

| Type | Kapan |
|---|---|
| `feat:` | Fitur baru |
| `fix:` | Bug fix |
| `chore:` | Tugas maintenance (update config, dependency, dll) |
| `refactor:` | Refactor kode tanpa perubahan behavior |
| `docs:` | Dokumentasi |
| `test:` | Tambah/update test |
| `style:` | Formatting, missing semicolon, dll (bukan logic) |

Contoh:
```
feat: add user registration endpoint
fix: handle empty request body on login
chore: update go version to 1.25.0
```

---

## Environment & Config

- Semua config lewat **environment variable**.
- Prefix env `APP_` untuk namespace.
- Contoh:
  - `APP_PORT=8080`
  - `APP_DB_URL=postgres://...`
  - `APP_JWT_SECRET=...`
- Jangan commit `.env` — tambahkan ke `.gitignore` dan buat `.env.example`.

---

## Mental Model

> "Write Go code the way the Go authors intended — simple, explicit, and compositional."

Prioritas:
1. **Readability** — kode dibaca lebih sering daripada ditulis.
2. **Simplicity** — solusi paling sederhana yang works.
3. **Correctness** — handle semua edge case dan error.
4. **Performance** — optimasi hanya jika terbukti diperlukan (profil dulu).
