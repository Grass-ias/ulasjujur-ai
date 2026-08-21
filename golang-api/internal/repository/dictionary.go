package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// IssueKeywordsCache sekarang berupa Map: map[keyword]category
var IssueKeywordsCache map[string]string

func InitializeDatabase(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("gagal membuka koneksi database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("database tidak merespons: %v", err)
	}

	log.Println("Database terkoneksi. Memulai pemuatan kamus kata kunci...")

	// Sesuaikan dengan nama tabel dan kolom dari temanmu
	rows, err := db.Query("SELECT keyword, category FROM issue_keywords")
	if err != nil {
		return fmt.Errorf("gagal menjalankan query kamus: %v", err)
	}
	defer rows.Close()

	// Inisialisasi map
	IssueKeywordsCache = make(map[string]string)
	
	count := 0
	for rows.Next() {
		var keyword, category string
		if err := rows.Scan(&keyword, &category); err != nil {
			log.Printf("Gagal membaca baris kamus: %v", err)
			continue
		}
		IssueKeywordsCache[keyword] = category
		count++
	}

	if count == 0 {
		return fmt.Errorf("tabel kamus kosong, pastikan script init SQL sudah dijalankan")
	}

	log.Printf("Berhasil memuat %d kata kunci ke dalam memory cache.", count)
	return nil
}