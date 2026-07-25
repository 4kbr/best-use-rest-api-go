# wrk — HTTP Benchmark Tool

## Ringkasan
- wrk = command line tool untuk benchmark HTTP server
- Bisa simulasi banyak koneksi concurrent
- Output: latency, request/detik, transfer rate
- Berguna untuk ngukur performa endpoint sebelum & sesudah optimasi

## Tugas

### 1. Cek pemahaman
- Apa yang diukur oleh wrk?
- Kenapa kita perlu benchmark?
- Apa bedanya `--threads` dan `--connections`?

### 2. Benchmark server.go
Jalankan server, lalu:
```bash
wrk -t4 -c100 -d10s https://localhost:3000/orders
```
Arti flag:
- `-t4` = 4 threads
- `-c100` = 100 koneksi concurrent
- `-d10s` = selama 10 detik

Catat output-nya: Latency, Req/Sec, Transfer/Sec.

### 3. Bandingkan
Coba dengan threads dan koneksi berbeda:
```bash
wrk -t1 -c10 -d10s https://localhost:3000/orders
wrk -t8 -c200 -d10s https://localhost:3000/orders
```
Apa perbedaan hasilnya? Kenapa?

### 4. Skenario
Tambah endpoint baru yang simulate delay (misal `time.Sleep(100 * time.Millisecond)`), lalu benchmark lagi. Berapa turun performanya?

## Tips
- Benchmark paling berguna untuk **before vs after** — bukan angka absolut
- Pastikan server jalan di terminal lain sebelum jalankan wrk
- Jangan benchmark di laptop yang sama dengan server untuk hasil akurat
- Kalau wrk error `Connection refused` → server belum jalan

## Kalau Masih Bingung
- Baca ulang: `docs/explanations/installing-wrk.md`

---

## Jawaban

### 1. Cek pemahaman

**Apa yang diukur oleh wrk?**
- **Latency**: rata-rata waktu response (ms)
- **Req/Sec**: jumlah request per detik
- **Transfer**: kecepatan data (MB/s)

**Kenapa kita perlu benchmark?**
Untuk tahu performa endpoint secara objektif — bukan "rasanya cepat" tapi ada angka. Berguna pas sebelum/ sesudah optimasi, atau bandingin 2 pendekatan.

**Apa bedanya `--threads` dan `--connections`?**
- `--threads` = jumlah CPU thread yang dipake wrk untuk generate request
- `--connections` = jumlah koneksi simultan ke server (yang stay open)
- Connections didistribusi ke threads. Misal `-t2 -c100` = tiap thread pegang 50 koneksi.

### 2. Benchmark server.go

Output kurang lebih (angka bisa beda):
```
Running 10s test @ https://localhost:3000/orders
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.24ms  492.20us  12.56ms   80.45%
    Req/Sec    20.15k     1.42k   23.78k    70.00%
  803456 requests in 10.01s, 125.98MB read
Requests/sec:  80345.67
Transfer/sec:      12.59MB
```

### 3. Bandingkan
- `-t1 -c10`: Latency lebih rendah (antrian pendek), tapi Req/Sec juga lebih rendah (hanya 1 thread).
- `-t8 -c200`: Latency naik (antrian panjang, contention), tapi Req/Sec mungkin naik sampai titik jenuh CPU.
- Biasanya ada **diminishing returns** — tambah thread/koneksi terus tidak selalu naik performa.

### 4. Skenario
Dengan `time.Sleep(100 * time.Millisecond)`:
```
Requests/sec: ~10 (turun drastis dari ~80000)
Latency: ~100ms
```
Karena setiap request nunggu 100ms — server cuma bisa handle ~10 request/detik per koneksi. Ini contoh kenapa **I/O blocking** itu musuh performa.
