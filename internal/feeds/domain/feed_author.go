package domain

import (
	shared_domain "daily-trends/go/internal/shared/domain"
)

type FeedAuthor string

func NewFeedAuthor(value string) (FeedAuthor, error) {
	if value == "" {
		return "", shared_domain.NewInvalidArgumentError("feed author must not be empty")
	}

	return FeedAuthor(value), nil
}
