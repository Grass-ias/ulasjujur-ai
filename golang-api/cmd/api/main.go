package main

import (
	"fmt"
	"log"
	"net/http"
	"os" // Tambahkan os untuk membaca environment variable

	"ulasjujur-api/internal/handler"
	"ulasjujur-api/internal/repository"
)

func main() {
	// Ambil DSN dari environment Docker, atau gunakan default localhost jika dijalankan manual
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/ulasjujur_db"
	}
	
	err := repository.InitializeDatabase(dsn)
	if err != nil {
		log.Printf("Peringatan: Gagal memuat database: %v. Pencocokan keyword mungkin tidak berfungsi penuh.", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("UlasJujur Golang API Service (v3) - Running"))
	})

	http.HandleFunc("/api/upload", handler.UploadCSVHandler)

	port := "8080" // Port internal di dalam container
	fmt.Printf("🚀 Server Golang berjalan di port %s...\n", port)
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}