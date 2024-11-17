package application

import (
	"context"
	"daily-trends/go/internal/feeds/domain"
)

type FeedSearcher struct {
	repository domain.FeedRepository
}

func NewFeedSearcher(repository domain.FeedRepository) *FeedSearcher {
	return &FeedSearcher{repository}
}

func (ff *FeedSearcher) Execute() ([]*domain.Feed, error) {
	ctx := context.Background()

	feeds, err := ff.repository.Search(ctx)
	if err != nil {
		return nil, err
	}

	return feeds, nil
}
