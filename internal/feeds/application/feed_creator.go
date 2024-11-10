package application

import (
	"context"
	"daily-trends/go/internal/feeds/domain"
)

type FeedCreator struct {
	repository domain.FeedRepository
}

func NewFeedCreator(repository domain.FeedRepository) *FeedCreator {
	return &FeedCreator{repository}
}

type FeedCreatorDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Source      string `json:"source"`
}

func (fc *FeedCreator) Execute(id string, dto *FeedCreatorDTO) error {
	ctx := context.Background()

	feed, err := domain.NewFeed(domain.WithId(id), domain.WithTitle(dto.Title), domain.WithAuthor(dto.Author), domain.WithSource(dto.Source), domain.WithDescription(dto.Description))
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
