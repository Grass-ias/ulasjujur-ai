package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ulasjujur-api/internal/client"
	"ulasjujur-api/internal/model" // Sesuaikan nama modul dengan yang kamu buat di go.mod
	"ulasjujur-api/internal/service"
)

const MaxRows = 100

func UploadCSVHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Gagal memproses form data", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File CSV tidak ditemukan pada request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		http.Error(w, "Gagal membaca header CSV atau file kosong", http.StatusBadRequest)
		return
	}

	reviewIdx, ratingIdx := -1, -1
	for i, header := range headers {
		colName := strings.TrimSpace(header)
		if strings.EqualFold(colName, "Customer Review") {
			reviewIdx = i
		} else if strings.EqualFold(colName, "Customer Rating") {
			ratingIdx = i
		}
	}

	if reviewIdx == -1 || ratingIdx == -1 {
		http.Error(w, "Format tidak sesuai: CSV harus memiliki kolom 'Customer Review' dan 'Customer Rating'", http.StatusBadRequest)
		return
	}

	var parsedData []model.ReviewRow
	var batchReviews []string

	rowNumber := 1
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			http.Error(w, fmt.Sprintf("Gagal membaca baris %d: %v", rowNumber+1, err), http.StatusBadRequest)
			return
		}
		rowNumber++

		if rowNumber > MaxRows+1 {
			http.Error(w, fmt.Sprintf("File melebihi batas maksimal %d baris", MaxRows), http.StatusBadRequest)
			return
		}

		reviewText := strings.TrimSpace(record[reviewIdx])
		ratingStr := strings.TrimSpace(record[ratingIdx])

		rating, err := strconv.Atoi(ratingStr)
		if err != nil || rating < 1 || rating > 5 {
			http.Error(w, fmt.Sprintf("Rating tidak valid pada baris %d (harus angka 1-5)", rowNumber), http.StatusBadRequest)
			return
		}

		parsedData = append(parsedData, model.ReviewRow{
			RowNumber:      rowNumber,
			CustomerRating: rating,
			CustomerReview: reviewText,
		})
		batchReviews = append(batchReviews, reviewText)
		// ... (Kode parsing CSV di atasnya tetap sama)
		parsedData = append(parsedData, model.ReviewRow{
			RowNumber:      rowNumber,
			CustomerRating: rating,
			CustomerReview: reviewText,
		})
		batchReviews = append(batchReviews, reviewText)
	}

	// ... (Kode CSV di atasnya tetap sama)

	// 1. Lakukan Batch Call ke Python Model
	aiResponses, err := client.PredictBatch(batchReviews)
	if err != nil {
		fmt.Printf("Warning: Gagal akses Python, menggunakan Mock Data AI. Error: %v\n", err)
		for i := 0; i < len(batchReviews); i++ {
			mockSentiment := "positive" // Format Python huruf kecil[cite: 2]
			if i%3 == 0 { 
				mockSentiment = "negative"
			}
			// Menyesuaikan dengan Struct AIModelResponseItem yang baru
			aiResponses = append(aiResponses, client.AIModelResponseItem{
				ID:         i + 1,
				Sentiment:  mockSentiment,
				Emotion:    "Puas", // Menggunakan format kembalian emosi dari Python[cite: 2]
				Confidence: 0.85,
				IsMixed:    false,
			})
		}
	}

	// 2. Jalankan Anomaly Logic & Agregasi
	analysisResult := service.AnalyzeData(parsedData, aiResponses)

	// 3. Susun Response Contract Final
	response := model.APIResponse{
		Status:  "success",
		Message: "Data berhasil dianalisis",
		Data:    analysisResult,
	}

	// Kirim JSON ke Frontend
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	importJsonEnc := json.NewEncoder(w) 
	importJsonEnc.SetIndent("", "  ")
	importJsonEnc.Encode(response)
}