package application_test

import (
	"context"
	"daily-trends/go/internal/feeds/application"
	application_mothers "daily-trends/go/internal/feeds/application/mothers"
	"daily-trends/go/internal/feeds/domain"
	"daily-trends/go/internal/feeds/domain/mocks"
	domain_mothers "daily-trends/go/internal/feeds/domain/mothers"
	shared_domain "daily-trends/go/internal/shared/domain"
	shared_mocks "daily-trends/go/internal/shared/domain/mocks"
	"errors"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeedCreator(t *testing.T) {
	t.Parallel()

	repo := mocks.NewFeedRepository(t)
	clock := shared_mocks.NewClockMock()
	creator := application.NewFeedCreator(repo, clock)

	t.Run("should return an InvalidArgumentError if the feed could not be instantiated", func(t *testing.T) {
		id := shared_domain.NewRandomUUID()
		creatorDTO := application_mothers.NewFeedCreatorInvalidDTO()

		err := creator.Execute(id.Value(), &creatorDTO)

		assert.ErrorAs(t, err, &shared_domain.InvalidArgumentError{})
	})

	t.Run("should return an error if the search of the feed fails", func(t *testing.T) {
		ctx := context.Background()

		id := shared_domain.NewRandomUUID()
		feedId, _ := domain.NewFeedId(id.Value())
		creatorDTO := application_mothers.NewFeedCreatorDTO()
		errorReturned := errors.New("unknown error")

		repo.On("FindById", ctx, *feedId).Return(nil, errorReturned)

		err := creator.Execute(id.Value(), &creatorDTO)

		assert.EqualError(t, err, errorReturned.Error())
	})

	t.Run("should return a FeedAlreadyExistsError if the feed already exists", func(t *testing.T) {
		ctx := context.Background()

		id := shared_domain.NewRandomUUID()
		feedId, _ := domain.NewFeedId(id.Value())
		creatorDTO := application_mothers.NewFeedCreatorDTO()
		savedFeed, _ := domain_mothers.NewRandomFeed(clock)

		repo.On("FindById", ctx, *feedId).Return(savedFeed, nil)

		err := creator.Execute(id.Value(), &creatorDTO)

		assert.ErrorAs(t, err, &domain.FeedAlreadyExistsError{})
	})

	t.Run("should return an error if the save of the feed fails", func(t *testing.T) {
		ctx := context.Background()

		id := shared_domain.NewRandomUUID()
		feedId, _ := domain.NewFeedId(id.Value())
		creatorDTO := application_mothers.NewFeedCreatorDTO()
		feedToSave, _ := domain_mothers.NewFeedFromProps(clock, id.Value(), creatorDTO)
		errorReturned := errors.New("unknown error")

		repo.On("FindById", ctx, *feedId).Return(nil, nil)
		repo.On("Save", ctx, feedToSave).Return(errorReturned)

		err := creator.Execute(id.Value(), &creatorDTO)

		assert.EqualError(t, err, errorReturned.Error())
		repo.AssertExpectations(t)
	})

	t.Run("should return nil if the feed is saved successfully", func(t *testing.T) {
		ctx := context.Background()

		id := shared_domain.NewRandomUUID()
		feedId, _ := domain.NewFeedId(id.Value())
		creatorDTO := application_mothers.NewFeedCreatorDTO()
		feedToSave, _ := domain_mothers.NewFeedFromProps(clock, id.Value(), creatorDTO)

		repo.On("FindById", ctx, *feedId).Return(nil, nil)
		repo.On("Save", ctx, feedToSave).Return(nil)

		err := creator.Execute(id.Value(), &creatorDTO)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})
}
