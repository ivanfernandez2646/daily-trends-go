package domain

import (
	shared_domain "daily-trends/go/internal/shared/domain"
	"net/url"
)

type FeedUrl *string

func NewFeedUrl(value string) (FeedUrl, error) {
	if value == "" {
		return nil, nil
	}

	parsedURL, err := url.Parse(value)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {

		return nil, shared_domain.NewInvalidArgumentError("feed url must be valid")
	}

	val := parsedURL.String()
	return FeedUrl(&val), nil
}
