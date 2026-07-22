# Port

## Apa itu Port?

Port adalah nomor "pintu" yang dipakai komputer untuk membedakan satu service dengan service lain. Satu komputer cuma punya 1 koneksi internet, tapi bisa menjalankan banyak program yang butuh internet. Port solusinya.

---

## Analogi: Gedung Apartemen

Bayangkan komputer adalah **sebuah gedung apartemen besar**.

| Konsep                 | Analogi                                        |
| ---------------------- | ---------------------------------------------- |
| **IP Address**         | Alamat gedung (misal: Jl. Merdeka No. 10)      |
| **Port**               | Nomor unit apartemen (Unit 101, Unit 202, dll) |
| **Pak Pos (Internet)** | Datang ke alamat gedung, cari unit tertentu    |
| **Service (Server)**   | Penghuni setiap unit                           |

**Cara kerja:**

1. Pak pos datang ke alamat gedung (`127.0.0.1`).
2. Dia lihat nomor unit yang dituju (`:8080`).
3. Dia ketuk pintu unit 8080.
4. Penghuni unit 8080 (server Go kita) buka pintu dan kasih response.

**Kenapa perlu unit?** Karena dalam satu gedung ada banyak penghuni. Kalau pak pos teriak "HALOO!" begitu sampai gedung, nggak jelas siapa yang dituju. Tapi kalau dia bilang "Pak pos untuk Unit 8080", jelas.

---

## Kenapa Port Diperlukan?

Satu komputer bisa menjalankan banyak service sekaligus:

- Browser lagi buka 10 tab
- Aplikasi chat (WhatsApp/Telegram)
- Game online
- Server Go kita
- Spotify streaming musik

Semua butuh internet. Tanpa port, bakal kacau — data Spotify nyasar ke server Go. Port bikin setiap aplikasi punya "antrian" sendiri.

---

## Format `:8080` di Go

```
":8080"
```

Kenapa ada `:` di depan? Ini format **`host:port`**.

| Format               | Arti                                                                  |
| -------------------- | --------------------------------------------------------------------- |
| `":8080"`            | Host dikosongin = dengerin dari semua jaringan (localhost, wifi, dll) |
| `"localhost:8080"`   | Cuma dengerin dari komputer ini aja                                   |
| `"192.168.1.5:8080"` | Cuma dengerin dari IP tertentu                                        |

Jadi `:8080` singkatan dari "dengerin port 8080 dari mana aja".

---

## Port yang Sering Dipakai

| Port    | Dipakai Oleh                               |
| ------- | ------------------------------------------ |
| `80`    | HTTP (website biasa)                       |
| `443`   | HTTPS (website aman)                       |
| `5432`  | PostgreSQL (database)                      |
| `3306`  | MySQL / MariaDB                            |
| `6379`  | Redis (cache)                              |
| `27017` | MongoDB                                    |
| `8080`  | Development server (biasanya pengganti 80) |
| `3000`  | React / Node.js dev server                 |

---

## Yang Perlu Diingat

1. **Port < 1024 butuh akses admin (root).** Soalnya ini port standar yang dipake service penting (80, 443). Makanya pas开发 kita pakai `:8080`.
2. **Dua program gabisa pakai port yang sama.** Kalau port 8080 udah kepake, kita harus pilih port lain.
3. **Port cuma angka.** Nggak ada bedanya port 80 sama 8080 secara teknis — cuma beda konvensi aja.
