package dto

import "time"

type CreateShortLinkRequest struct {
	URL string `json:"url" binding:"required,url"`
	ShortCode string `json:"short_code,omitempty" binding:"omitempty,min=3,max=32"`
}

type CreateShortLinkResponse struct {
	ShortCode   string    `json:"short_code"`
	ShortURL    string    `json:"short_url"`
	OriginalUrl string    `json:"original_url"`
	UserID      uint64    `json:"user_id"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}
