package domain

import (
	shared_domain "daily-trends/go/internal/shared/domain"
)

type FeedTitle string

func NewFeedTitle(value string) (*FeedTitle, error) {
	if value == "" {
		return nil, shared_domain.NewInvalidArgumentError("feed title must not be empty")
	}

	val := FeedTitle(value)
	return &val, nil
}
