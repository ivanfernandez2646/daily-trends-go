package domain

import "context"

type FeedRepository interface {
	Save(ctx context.Context, feed *Feed) error
	FindById(ctx context.Context, id FeedId) (*Feed, error)

	// TODO: add criteria
	Search(ctx context.Context) ([]*Feed, error)
}
