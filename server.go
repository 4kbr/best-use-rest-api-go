// package main adalah package spesial yang menghasilkan binary executable.
// Setiap program Go yang bisa dijalankan harus punya func main() di package main.
package main

// import digunakan untuk mengambil package lain.
// fmt: format teks (print ke console atau response).
// net/http: package untuk web server, request, response.
// crypto/tls: package untuk TLS (encryption, certificate, handshake).
// Dipakai untuk konfigurasi keamanan HTTPS.
import (
	"crypto/tls"

	"fmt"
	"log"
	"net/http"

	// golang.org/x/net/http2: package external dari Go team.
	// http2.ConfigureServer() enable HTTP/2 di atas koneksi TLS.
	// Ini adalah "golang.org/x/net" — extension package di luar std lib.
	"golang.org/x/net/http2"
)

// func main adalah entry point. Semua program mulai dari sini.
func main() {
	// http.HandleFunc mendaftarkan handler untuk suatu path.
	// Setiap kali ada request ke "/orders", fungsi ini akan dijalankan.
	// w = tempat menulis response, r = data request dari client.
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {

		// logRequestDetails: cetak HTTP version & TLS version ke terminal.
		// Berguna pas debugging — tahu client pake HTTP/1.1 atau HTTP/2,
		// dan apakah koneksi pakai TLS atau tidak.
		logRequestDetails(r)

		// fmt.Fprintf menulis teks yang diformat ke ResponseWriter.
		// Teks ini akan dikirim sebagai response ke client.
		fmt.Fprintf(w, "handling incoming orders")
	})

	// http.HandleFunc mendaftarkan handler ke path "/users".
	// Pola yang sama seperti "/orders", beda path dan response.
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

		// logRequestDetails: cetak HTTP version & TLS version ke terminal.
		// Berguna pas debugging — tahu client pake HTTP/1.1 atau HTTP/2,
		// dan apakah koneksi pakai TLS atau tidak.
		logRequestDetails(r)

		// fmt.Fprintf menuliskan response string ke client.
		fmt.Fprintf(w, "handling incoming users")
	})

	// deklarasi variable port dengan tipe yang di-infer oleh Go.
	// ":=" adalah short variable declaration (deklarasi + assignment).
	port := 3000

	// cert.pem = file certificate (public, self-signed).
	// key.pem  = file private key (RAHASIA, jangan pernah di-commit).
	// Keduanya dihasilkan dari: openssl req -x509 -newkey rsa:2048 ...
	cert := "cert.pem"
	key := "key.pem"

	// tls.Config: konfigurasi TLS untuk server.
	// MinVersion: batasi protokol TLS minimal versi 1.2.
	// TLS 1.0 dan 1.1 sudah deprecated — rawan serangan.
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// &http.Server{}: kita bikin instance server sendiri (bukan pake default).
	// Kenapa? Karena kita perlu ngasih TLSConfig dan nanti enable HTTP/2.
	// Handler nil: pake DefaultServeMux yang udah kita daftarin HandleFunc.
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   nil,
		TLSConfig: tlsConfig,
	}

	// http2.ConfigureServer: enable HTTP/2 di server ini.
	// HTTP/2 butuh TLS — browser cuma support h2 (HTTP/2 via TLS),
	// bukan h2c (HTTP/2 cleartext / tanpa TLS).
	http2.ConfigureServer(server, &http2.Server{})

	fmt.Println("server is running on port:", port)

	// ListenAndServeTLS: sama seperti ListenAndServe, tapi dengan TLS.
	// Parameter 1: path ke certificate file (cert.pem).
	// Parameter 2: path ke private key file (key.pem).
	// Fungsi ini blocking — jalan terus sampai server dimatikan.
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("could not start server", err)
	}

	// HTTP 1.1 server without TLS
	// // http.ListenAndServe menghidupkan HTTP server.
	// // Parameter 1: alamat "host:port" (":3000").
	// // Parameter 2: handler (nil = pakai DefaultServeMux).
	// // Fungsi ini blocking — server jalan terus sampai dimatikan.
	// // tambah handle error
	// err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	//
	//	if err != nil {
	//		log.Fatalln("could not start server")
	//	}
}

// logRequestDetails: mencetak informasi HTTP request ke terminal.
// r.Proto = HTTP version dari request ("HTTP/1.1", "HTTP/2", dll).
// r.TLS   = nil kalau koneksi HTTP biasa, berisi struct TLS kalau HTTPS.
// Fungsi ini dipanggil di setiap handler buat debugging.
func logRequestDetails(r *http.Request) {
	// r.Proto: string yang berisi HTTP version.
	// Contoh: "HTTP/1.1", "HTTP/2", "HTTP/3".
	httpVersion := r.Proto
	fmt.Println("received request with HTTP version:", httpVersion)

	// r.TLS: kalau nil → koneksi HTTP biasa (tanpa enkripsi).
	// Kalau tidak nil → koneksi HTTPS (TLS aktif).
	// r.TLS.Version: uint16 yang menandakan versi TLS.
	if r.TLS != nil {
		tlsVersion := getTLSVersionName(r.TLS.Version)
		fmt.Println("received request with TLS version:", tlsVersion)
	} else {
		fmt.Println("received request without TLS")
	}
}

// getTLSVersionName: konversi angka versi TLS ke string yang mudah dibaca.
// tls.VersionTLSxx adalah constant dari package crypto/tls.
// Parameter: version uint16 dari r.TLS.Version.
// Return: string seperti "TLS 1.2", "TLS 1.3", dll.
func getTLSVersionName(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("default: versi %d", version)
	}
}
