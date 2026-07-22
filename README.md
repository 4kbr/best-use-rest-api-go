# best-practice-api

REST API sederhana sebagai learning ground untuk best practice Go. Stack minimal — pure std library `net/http`, tanpa framework eksternal.

## Tech Stack

- **Go** 1.25.0
- **Router** `net/http` (std lib)
- **No external framework** — semua dari standard library

## Cara Jalankan

```bash
go run server.go
```

Atau compile dulu:

```bash
go build -o server server.go && ./server
```

Server akan listen di `:8080`.

## Struktur Folder

```
.
├── cmd/
│   └── server/          # entry point (coming soon)
├── docs/
│   └── explanations/    # penjelasan konsep (port, dll)
├── internal/
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # middleware (logging, auth, CORS)
│   ├── model/           # struct/type definitions
│   ├── repository/      # data access layer
│   └── service/         # business logic
├── pkg/                 # public (exported) package
├── .gitignore
├── AGENTS.md            # panduan untuk AI agent
├── go.mod
├── go.sum
└── server.go            # entry point (akan pindah ke cmd/server/)
```

## Testing

```bash
go test ./...
```

Semua test pakai package `testing` bawaan Go, dengan table-driven test.

## Endpoint

| Method      | Path        | Status |
| ----------- | ----------- | ------ |
| (belum ada) | (belum ada) | 🚧 WIP |

Endpoint akan bertambah seiring development.

## Tujuan Project

- Belajar idiomatic Go — simple, explicit, compositional
- Penerapan clean architecture di Go (`internal/` sebagai private package)
- REST API patterns: method-based routing, JSON envelope response, consistent error codes
- Best practice testing (table-driven test, interface + mock)

## License

MIT
