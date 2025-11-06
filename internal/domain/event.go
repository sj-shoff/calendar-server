package domain

// Event представляет событие в календаре
type Event struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Date   string `json:"date"` // YYYY-MM-DD
	Title  string `json:"title"`
}
