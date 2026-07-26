# catatan (terlalu malas untuk ngoding jadi nyimak doang)

- sebelum: postman coba settings postman `Enable SSL certificate verification: on` hasilnya gagal
- buka settings postman, add certificate upload cert.pem dan key.pem, lalu coba get lagi dan masih gagal
- buat openssl.cnf, hapus cert.pem dan key.pem saat ini
- jalankan command `openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config openssl.cnf` otomatis membuat cert.pem dan key.pem
- go back to postman, settings, upload cert.pem di CA certificates
- coba send request lagi dan berhasil!!

- mTLS sangat spesifik use dan tidak umum digunakan, dibicarakan hanya untuk menambah pengetahuan bukan untuk diimplementasi
