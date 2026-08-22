package service

import (
	"strings"

	"ulasjujur-api/internal/client"
	"ulasjujur-api/internal/model"
	"ulasjujur-api/internal/repository"
)

// AnalyzeData menjalankan Anomaly Logic dan Keyword Integration
func AnalyzeData(rows []model.ReviewRow, aiResults []client.AIModelResponseItem) model.AnalysisResponseData {
	var totalPositive, totalNegative, totalNeutral int
	issueCounts := make(map[string]int)
	emotionCounts := make(map[string]int)

	mismatches := make([]model.MismatchExample, 0)

	for i, row := range rows {
		if i >= len(aiResults) {
			break
		}
		aiRes := aiResults[i]

		// 1. Agregasi Persentase Sentimen
		// Menggunakan strings.EqualFold karena Python mengirimkan "positive" / "negative" huruf kecil[cite: 2]
		if strings.EqualFold(aiRes.Sentiment, "positive") {
			totalPositive++
		} else if strings.EqualFold(aiRes.Sentiment, "negative") {
			totalNegative++
		} else {
			totalNeutral++
		}

		// 2. Agregasi Emosi Dominan
		emotionCounts[aiRes.Emotion]++

		// 3 & 4. KEYWORD INTEGRATION & ANOMALY DETECTION
        lowerReview := strings.ToLower(row.CustomerReview)
        foundIssue := false

        // Scan semua keyword di ulasan tanpa mempedulikan sentimen awal AI
        for keyword, category := range repository.IssueKeywordsCache {
            if strings.Contains(lowerReview, keyword) {
                foundIssue = true
                // Tambahkan hitungan metrik masalah kalau sentimen aslinya memang negatif
                if strings.EqualFold(aiRes.Sentiment, "negative") {
                    issueCounts[category]++
                }
            }
        }

        // LOGIKA ANOMALI BARU (MIXED SENTIMENT):
        // Jika pelanggan kasih rating bagus (>=4) ATAU AI nebaknya "Positive",
        // TAPI sistem kamus nemuin kata keluhan (foundIssue == true) itu berarti Mixed Sentiment
        if (row.CustomerRating >= 4 || strings.EqualFold(aiRes.Sentiment, "positive")) && foundIssue {
            mismatches = append(mismatches, model.MismatchExample{
                RowNumber:         row.RowNumber,
                CustomerRating:    row.CustomerRating,
                CustomerReview:    row.CustomerReview,
                AISentiment:       aiRes.Sentiment,
                AIDominantEmotion: aiRes.Emotion,
            })
        }
	}

	topIssues := make([]model.Issue, 0)
	for kw, count := range issueCounts {
		topIssues = append(topIssues, model.Issue{Keyword: kw, Count: count})
	}

	var maxEmotion string
	var maxCount int
	for emo, count := range emotionCounts {
		if count > maxCount {
			maxCount = count
			maxEmotion = emo
		}
	}

	calcPercent := func(val, total int) int {
		if total == 0 {
			return 0
		}
		return (val * 100) / total
	}

	totalReviews := len(rows)

	return model.AnalysisResponseData{
		Summary: model.Summary{
			TotalReviews: totalReviews,
			SentimentPercentage: map[string]int{
				"positive": calcPercent(totalPositive, totalReviews),
				"negative": calcPercent(totalNegative, totalReviews),
				"neutral":  calcPercent(totalNeutral, totalReviews),
			},
			DominantEmotion: maxEmotion,
		},
		TopIssues:        topIssues,
		MismatchExamples: mismatches,
	}
}