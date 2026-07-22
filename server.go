// package main adalah package spesial yang menghasilkan binary executable.
// Setiap program Go yang bisa dijalankan harus punya func main() di package main.
package main

// import digunakan untuk mengambil package lain.
// fmt: format teks (print ke console atau response).
// net/http: package untuk web server, request, response.
import (
	"fmt"
	"net/http"
)

// func main adalah entry point. Semua program mulai dari sini.
func main() {
	// http.HandleFunc mendaftarkan handler untuk suatu path.
	// Setiap kali ada request ke "/orders", fungsi ini akan dijalankan.
	// w = tempat menulis response, r = data request dari client.
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf menulis teks yang diformat ke ResponseWriter.
		// Teks ini akan dikirim sebagai response ke client.
		fmt.Fprintf(w, "handling incoming orders")
	})

	// http.HandleFunc mendaftarkan handler ke path "/users".
	// Pola yang sama seperti "/orders", beda path dan response.
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf menuliskan response string ke client.
		fmt.Fprintf(w, "handling incoming users")
	})

	// deklarasi variable port dengan tipe yang di-infer oleh Go.
	// ":=" adalah short variable declaration (deklarasi + assignment).
	port := 3000

	// fmt.Println mencetak log ke terminal untuk memberi tahu developer.
	fmt.Println("server is running on port:", port)

	// http.ListenAndServe menghidupkan HTTP server.
	// Parameter 1: alamat "host:port" (":3000").
	// Parameter 2: handler (nil = pakai DefaultServeMux).
	// Fungsi ini blocking — server jalan terus sampai dimatikan.
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}
