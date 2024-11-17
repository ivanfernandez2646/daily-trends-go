package application_test

import (
	"context"
	"daily-trends/go/internal/feeds/application"
	"daily-trends/go/internal/feeds/domain"
	"daily-trends/go/internal/feeds/domain/mocks"
	"daily-trends/go/internal/feeds/domain/mothers"
	shared_domain "daily-trends/go/internal/shared/domain"
	shared_mocks "daily-trends/go/internal/shared/domain/mocks"

	"errors"

	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeedFinder(t *testing.T) {
	t.Parallel()

	repo := mocks.NewFeedRepository(t)
	clock := shared_mocks.NewClockMock()
	finder := application.NewFeedFinder(repo)

	t.Run("should return an InvalidArgumentError when uuid is not valid", func(t *testing.T) {
		id := "wrong1234"

		feed, err := finder.Execute(id)

		assert.Nil(t, feed)
		assert.EqualError(t, err, fmt.Sprintf("invalid argument error: cannot parse uuid %s", id))
	})

	t.Run("should return an error if FindById fails", func(t *testing.T) {
		ctx := context.Background()

		id := shared_domain.NewRandomUUID()
		idValue := id.Value()
		errorReturned := errors.New("unknown error")

		feedId, _ := domain.NewFeedId(idValue)
		repo.On("FindById", ctx, *feedId).Return(nil, errorReturned)

		feed, err := finder.Execute(idValue)

		assert.Nil(t, feed)
		assert.EqualError(t, err, errorReturned.Error())
	})

	t.Run("should return a FeedNotFoundError if feed is not found", func(t *testing.T) {
		ctx := context.Background()

		id := shared_domain.NewRandomUUID()
		idValue := id.Value()
		feedId, _ := domain.NewFeedId(idValue)
		notFoundError := domain.NewFeedNotFoundError(*feedId)

		repo.On("FindById", ctx, *feedId).Return(nil, notFoundError)

		res, err := finder.Execute(feedId.Value())

		assert.Nil(t, res)
		assert.EqualError(t, err, notFoundError.Error())
	})

	t.Run("should return a Feed is found", func(t *testing.T) {
		ctx := context.Background()

		feed, _ := mothers.NewRandomFeed(clock)

		repo.On("FindById", ctx, feed.Id()).Return(feed, nil)

		res, err := finder.Execute(feed.Id().Value())

		assert.Nil(t, err)
		assert.Equal(t, *feed, *res)
	})
}
