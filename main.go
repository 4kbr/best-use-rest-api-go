// package main = program ini bisa dijalankan langsung
package main

import (
	// encoding/json = buat manipulasi JSON (marshal, unmarshal, decoder, encoder)
	"encoding/json"
	// fmt = buat print teks ke terminal
	"fmt"
	// os = akses sistem operasi (stdout, stdin, stderr, dll)
	"os"
	// strings = manipulasi teks (NewReader, Contains, Split, dll)
	"strings"
)

// tipe data struct (mirip class di bahasa lain, tapi tanpa method ribet)
// field Name dan Email pake tag `json:"..."`
// tag ini memberitahu Go: "kalau diubah ke JSON, panggil field-nya 'name' dan 'email'"
// kalau tanpa tag, Go pake nama field (Name, Email) — huruf besar di depan
// huruf kecil untuk field JSON itu konvensi umum di REST API
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// func main = titik awal eksekusi program
func main() {
	// ============================================================
	// MARSHAL — ngubah struct Go jadi JSON (encode)
	// ============================================================

	// bikin struct User dengan dua field
	user := User{Name: "alice", Email: "alice@example.com"}

	// json.Marshal(user): ubah struct user ke format JSON
	// return-nya 2 nilai:
	//   1. []byte (kumpulan byte / data mentah) — JSON result-nya
	//   2. error — kalau gagal (misal: field-nya nggak bisa di-JSON-kan)
	jsonData, err := json.Marshal(user)

	// error handling: WAJIB di Go
	// kalau err bukan nil, cetak errornya dan stop program (pakai return)
	if err != nil {
		fmt.Println("error marshalling JSON:", err)
		return
	}

	// jsonData tipenya []byte — kita ubah ke string biar bisa dibaca
	// %s = format untuk string, bisa juga pake string(jsonData)
	fmt.Printf("JSON result: %s\n", jsonData)

	// ============================================================
	// UNMARSHAL — ngubah JSON string jadi struct Go (decode)
	// ============================================================

	// JSON string sebagai data input kita
	// pake backtick ` ` biar string bisa multi-line tanpa escape character
	jsonString := `{"name":"bob","email":"bob@example.com"}`

	// siapin variable tempat hasil unmarshal
	var decodedUser User

	// json.Unmarshal: ubah JSON string ke struct
	// parameter 1: []byte (JSON mentah) — kita konversi dari string pake []byte(...)
	// parameter 2: pointer ke struct (&decodedUser)
	// kenapa pake pointer (&)? karena Go pass by value — kalau tanpa pointer,
	// decodedUser cuma di-copy, isi aslinya nggak berubah
	err = json.Unmarshal([]byte(jsonString), &decodedUser)
	if err != nil {
		fmt.Println("error unmarshalling JSON:", err)
		return
	}

	// cetak hasil struct setelah diisi dari JSON
	fmt.Printf("Decoded struct: Name=%s, Email=%s\n", decodedUser.Name, decodedUser.Email)

	// ============================================================
	// MARSHAL INDENT — JSON dengan formatting/indentasi (biar rapi)
	// ============================================================

	// json.MarshalIndent mirip Marshal biasa, tapi ada parameter tambahan:
	//   prefix: string di awal tiap baris (biasanya "")
	//   indent: string buat indentasi tiap level (biasanya "  " atau "\t")
	prettyJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Println("error marshalling pretty JSON:", err)
		return
	}

	// cetak JSON dengan indentasi — lebih enak dibaca manusia
	fmt.Println("Pretty JSON:")
	fmt.Println(string(prettyJSON))

	// ============================================================
	// NEWDECODER — baca JSON dari io.Reader (streaming decode)
	// ============================================================

	// strings.NewReader: bungkus string jadi io.Reader
	// io.Reader adalah interface yang bisa dibaca byte-by-byte
	// Kenapa perlu Reader? Karena data JSON bisa datang dari mana aja:
	//   - file (os.File)
	//   - network (http.Response.Body)
	//   - keyboard (os.Stdin)
	//   - string (strings.NewReader)
	reader := strings.NewReader(`{"name":"charlie","email":"charlie@example.com"}`)

	// json.NewDecoder: bikin decoder yang baca dari reader
	// Decoder membaca JSON secara STREAMING — satu per satu, nggak perlu loading semua data
	decoder := json.NewDecoder(reader)

	// siapin variable tempat hasil decode
	var streamUser User

	// decoder.Decode(&user): baca satu JSON object dari stream, lalu decode ke struct
	// parameter harus pointer (&streamUser) biar struct-nya diisi
	// Decode bisa dipanggil berkali-kali kalau stream-nya punya banyak JSON object
	err = decoder.Decode(&streamUser)
	if err != nil {
		fmt.Println("error decoding JSON:", err)
		return
	}

	fmt.Println("--- JSON DECODER (streaming) ---")
	fmt.Printf("Name: %s, Email: %s\n", streamUser.Name, streamUser.Email)

	// ============================================================
	// NEWENCODER — tulis JSON ke io.Writer (streaming encode)
	// ============================================================

	// os.Stdout: writer ke terminal (standar output)
	// json.NewEncoder: bikin encoder yang nulis ke writer
	// Encoder nulis JSON secara STREAMING — langsung dikirim ke writer, bukan disimpan di memory
	encoder := json.NewEncoder(os.Stdout)

	fmt.Println("--- JSON ENCODER (to stdout) ---")

	// encoder.Encode(user): ubah struct user ke JSON, lalu tulis ke writer (stdout)
	// Mirip json.Marshal, tapi outputnya langsung ke writer, bukan return []byte
	// Keuntungan: nggak perlu nyimpen JSON di memory — cocok buat data besar
	err = encoder.Encode(user)
	if err != nil {
		fmt.Println("error encoding JSON:", err)
		return
	}

	// ============================================================
	// KESIMPULAN — Kapan Pakai Apa?
	// ============================================================
	//
	// | Operasi            | API                   | Cocok Untuk                                  |
	// |--------------------|-----------------------|----------------------------------------------|
	// | Struct → JSON      | json.Marshal          | Data kecil, butuh []byte buat dikirim/disimpan |
	// | Struct → JSON rapi | json.MarshalIndent    | Debugging, config file, human-readable output |
	// | JSON → Struct      | json.Unmarshal        | JSON udah di memory ([]byte/string)           |
	// | JSON → Struct      | json.NewDecoder+Decode| JSON dari file/network/stream — data besar   |
	// | Struct → JSON      | json.NewEncoder+Encode| Langsung output ke file/network/stdout       |
	//
	// Aturan praktis:
	// - Data udah ada di variable ([]byte/string) → Marshal / Unmarshal
	// - Data datang dari file / koneksi / stream → NewDecoder
	// - Data mau ditulis langsung ke file / koneksi / stdout → NewEncoder
}
