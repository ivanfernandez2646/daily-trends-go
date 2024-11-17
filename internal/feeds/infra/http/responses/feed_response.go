package responses

import (
	"daily-trends/go/internal/feeds/domain"
	"time"
)

type GetFeedResponseDTO struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Author      string  `json:"author"`
	Source      string  `json:"source"`
	Url         *string `json:"url"`
	CreatedAt   string  `json:"createdAt"`
}

func NewFeedGetResponse(feed *domain.Feed) *GetFeedResponseDTO {
	id := feed.Id()
	feedDTO := &GetFeedResponseDTO{
		Id:          id.Value(),
		Title:       string(feed.Title()),
		Description: feed.Description(),
		Author:      string(feed.Author()),
		Source:      feed.Source().String(),
		Url:         feed.Url(),
		CreatedAt:   feed.CreatedAt().Format(time.RFC3339),
	}

	return feedDTO
}

func NewFeedHomeGetResponse(feeds []*domain.Feed) []*GetFeedResponseDTO {
	var res []*GetFeedResponseDTO

	for _, feed := range feeds {
		feedDTO := NewFeedGetResponse(feed)
		res = append(res, feedDTO)
	}

	return res
}
