package domain

import (
	"time"

	shared_domain "daily-trends/go/internal/shared/domain"
)

type Feed struct {
	id          FeedId
	title       FeedTitle
	description FeedDescription
	author      FeedAuthor
	source      FeedSource
	url         FeedUrl
	createdAt   time.Time
}

func NewFeed(clock shared_domain.Clock, options ...FeedOption) (*Feed, error) {
	feed := &Feed{
		createdAt: clock.Now(),
	}
	for _, option := range options {
		err := option(feed)
		if err != nil {
			return nil, err
		}
	}

	return feed, nil
}

func NewFeedOnlyFromOptions(options ...FeedOption) (*Feed, error) {
	feed := &Feed{}
	for _, option := range options {
		err := option(feed)
		if err != nil {
			return nil, err
		}
	}

	return feed, nil
}

func (f *Feed) Id() FeedId {
	return f.id
}

func (f *Feed) Title() FeedTitle {
	return f.title
}

func (f *Feed) Author() FeedAuthor {
	return f.author
}

func (f *Feed) Description() FeedDescription {
	if f.description == nil {
		return nil
	}

	copy := *f.description
	return &copy
}

func (f *Feed) Source() FeedSource {
	return f.source
}

func (f *Feed) Url() FeedUrl {
	if f.url == nil {
		return nil
	}

	url := *f.url
	return &url
}

func (f *Feed) CreatedAt() time.Time {
	return f.createdAt
}
