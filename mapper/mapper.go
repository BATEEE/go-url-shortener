package mapper

import (
	"simple-shortener/domain"
	"simple-shortener/domain/dto"
)

func ToCreateShortResponse(d domain.Link) dto.CreateShortLinkResponse {
	return dto.CreateShortLinkResponse{
		ShortCode:   d.ShortCode,
		OriginalUrl: d.OriginalUrl,
		UserID:      d.UserID,
		Clicks:      d.Clicks,
		CreatedAt:   d.CreatedAt,
	}
}
