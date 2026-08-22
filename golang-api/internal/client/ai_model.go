package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Menyesuaikan dengan ReviewInput di main.py
type AIModelRequestItem struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type AIModelBatchRequest struct {
	Reviews []AIModelRequestItem `json:"reviews"`
}

// Menyesuaikan dengan format kembalian di main.py
type AIModelResponseItem struct {
	ID         int     `json:"id"`
	Sentiment  string  `json:"sentiment"`
	Emotion    string  `json:"emotion"`
	Confidence float64 `json:"confidence"`
}

type AIModelBatchResponse struct {
	Results []AIModelResponseItem `json:"results"`
}

// PredictBatch memanggil main.py FastAPI
// PredictBatch memanggil main.py FastAPI
func PredictBatch(reviews []string) ([]AIModelResponseItem, error) {
	// Ambil URL dari Environment Variable Docker
	url := os.Getenv("PYTHON_API_URL")
	if url == "" {
		// Fallback jika dijalankan manual tanpa Docker
		url = "http://localhost:5000/predict-batch"
	}
	// Menyusun array object sesuai ekspektasi Python
	var requestItems []AIModelRequestItem
	for i, text := range reviews {
		requestItems = append(requestItems, AIModelRequestItem{
			ID:   i + 1, // Berikan ID sekuensial sederhana
			Text: text,
		})
	}

	payload := AIModelBatchRequest{Reviews: requestItems}
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan encode batch request: %v", err)
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("gagal memanggil Python service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("model AI merespons dengan status: %d", resp.StatusCode)
	}

	// Tangkap respons JSON dari Python
	var result AIModelBatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal memilah respons dari model AI: %v", err)
	}

	return result.Results, nil
}