// WAJIB: Setiap file Go punya package.
// "main" adalah package spesial — dia jadi entry point program.
// Kalau package-nya bukan "main", dia jadi library yang bisa di-import.
package main

// import: cara ambil package bawaan Go atau punya orang lain.
// - fmt: buat format teks (print ke console, tulis ke response, dll)
// - log: buat logging (catat error/info)
// - net/http: segalanya tentang web (server, request, response)
import (
	"fmt"
	"log"
	"net/http"
)

// func main() = fungsi pertama yang jalan saat program dijalankan.
// Ini CERITA UTAMA program kita. Semua dimulai dari sini.
func main() {
	// http.HandleFunc: "Hei Go, kalau ada yang request ke path "/",
	// jalankan fungsi ini."
	//
	// w http.ResponseWriter = tempat kita nulis jawaban.
	// ibaratnya kita punya kertas kosong, kita tulis "hello from server!"
	// lalu Go kirimkan kertas itu ke client (browser, curl, dll).
	//
	// r *http.Request = isi surat dari client.
	// Di sini ada method (GET/POST), URL, header, body, dll.
	// Nanti kita bakal sering lihat "r" ini.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// fmt.Fprintln: cetak teks ke "writer" (w).
		// F println = File + Print Line.
		// w adalah ResponseWriter — jadi teks ini masuk ke response.
		fmt.Fprintln(w, "hello from server!")
	})

	// const: bikin nilai tetap yang nggak bisa berubah.
	// port: nama variable-nya (lowercase karena cuma dipakai di sini).
	// string: tipe data-nya (kata/kalimat).
	// ":8080": nilai-nya. Artinya server dengerin di port 8080.
	// Kenapa ada ":" di depan? Itu aturan Go: port harus diawali ":".
	const port string = ":8080"

	// Alternatif: pake := (deklarasi + assign otomatis).
	// Go tebak tipe data dari nilai di kanan.
	// Baris ini sengaja di-comment (tidak aktif).
	// port := ":8080"

	// fmt.Println: cetak teks ke terminal (STDOUT).
	// Berguna buat ngasih tau status ke developer yang menjalankan server.
	fmt.Println("server listening on port:", port)

	// http.ListenAndServe: "Hidupkan server dan dengerin request!"
	// Parameter 1: port (":8080")
	// Parameter 2: handler (nil = pake default, yaitu HandleFunc di atas)
	//
	// Fungsi ini NGELOCK program — server jalan terus sampai dimatiin.
	// err := ... karena ListenAndServe MENGEMBALIKAN error kalau gagal.
	err := http.ListenAndServe(port, nil)

	// if err != nil = cek apakah ada error.
	// nil = kosong/tidak ada. Kalau err bukan nil, artinya ada masalah.
	// Ini WAJIB dilakukan di Go — kalau ada error, kita harus tangani.
	if err != nil {
		// log.Fatalln: cetak error ke terminal, lalu matikan program.
		// Fatal = berhenti total. Ini cuma dipakai di func main.
		// Di fungsi lain, kita return error, jangan pake Fatal.
		log.Fatalln("error when starting server", err)
	}
}
