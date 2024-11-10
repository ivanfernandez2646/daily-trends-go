package application

import (
	"context"
	"daily-trends/go/internal/feeds/domain"
)

type FeedFinder struct {
	repository domain.FeedRepository
}

func NewFeedFinder(repository domain.FeedRepository) *FeedFinder {
	return &FeedFinder{repository}
}

func (ff *FeedFinder) Execute(id string) (*domain.Feed, error) {
	ctx := context.Background()

	feedId, err := domain.NewFeedId(id)
	if err != nil {
		return nil, err
	}

	feed, err := ff.repository.FindById(ctx, *feedId)
	if err != nil {
		return nil, err
	}

	if feed == nil {
		return nil, domain.NewFeedNotFoundError(*feedId)
	}

	return feed, nil
}
