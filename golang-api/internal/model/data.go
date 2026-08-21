package model

// ReviewRow merepresentasikan satu baris data valid dari CSV mock data PRDECT-ID
type ReviewRow struct {
	RowNumber      int
	CustomerRating int
	CustomerReview string
}