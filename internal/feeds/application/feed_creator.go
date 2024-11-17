package application

import (
	"context"
	"daily-trends/go/internal/feeds/domain"
	shared_domain "daily-trends/go/internal/shared/domain"
)

type FeedCreator struct {
	repository domain.FeedRepository
	clock      shared_domain.Clock
}

func NewFeedCreator(repository domain.FeedRepository, clock shared_domain.Clock) *FeedCreator {
	return &FeedCreator{repository, clock}
}

type FeedCreatorDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Source      string `json:"source"`
}

func (fc *FeedCreator) Execute(id string, dto *FeedCreatorDTO) error {
	ctx := context.Background()

	feed, err := domain.NewFeed(fc.clock, domain.WithId(id), domain.WithTitle(dto.Title), domain.WithAuthor(dto.Author), domain.WithSource(dto.Source), domain.WithDescription(dto.Description))
	if err != nil {
		return err
	}

	exists, err := fc.repository.FindById(ctx, feed.Id())
	if err != nil {
		return err
	}

	if exists != nil {
		return domain.NewFeedAlreadyExistsError(feed.Id())
	}

	err = fc.repository.Save(ctx, feed)
	if err != nil {
		return err
	}

	return nil
}
