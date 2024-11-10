package domain

import (
	shared_domain "daily-trends/go/internal/shared/domain"
)

type FeedTitle string

func NewFeedTitle(value string) (FeedTitle, error) {
	if value == "" {
		return "", shared_domain.NewInvalidArgumentError("feed title must not be empty")
	}

	return FeedTitle(value), nil
}
