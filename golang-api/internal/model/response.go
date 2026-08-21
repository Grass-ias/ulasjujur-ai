package model

type Issue struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}

type MismatchExample struct {
	RowNumber         int    `json:"row_number"`
	CustomerRating    int    `json:"customer_rating"`
	CustomerReview    string `json:"customer_review"`
	AISentiment       string `json:"ai_sentiment"`
	AIDominantEmotion string `json:"ai_dominant_emotion"`
}

type Summary struct {
	TotalReviews        int            `json:"total_reviews"`
	SentimentPercentage map[string]int `json:"sentiment_percentage"`
	DominantEmotion     string         `json:"dominant_emotion"`
}

type AnalysisResponseData struct {
	Summary          Summary           `json:"summary"`
	TopIssues        []Issue           `json:"top_issues"`
	MismatchExamples []MismatchExample `json:"mismatch_examples"`
}

type APIResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Data    AnalysisResponseData `json:"data"`
}