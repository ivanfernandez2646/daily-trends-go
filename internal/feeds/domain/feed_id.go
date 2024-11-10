package domain

import shared_domain "daily-trends/go/internal/shared/domain"

type FeedId struct {
	shared_domain.UUID
}

func NewFeedId(value string) (*FeedId, error) {
	val, err := shared_domain.NewUUID(value)
	if err != nil {
		return nil, err
	}

	feedId := FeedId{*val}

	return &feedId, err
}
