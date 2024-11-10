package application

import (
	"context"
	"daily-trends/go/internal/feeds/domain"
)

type FeedScraperCreator struct {
	scraper    domain.FeedScraper
	extractors []domain.FeedContentExtractor
	repository domain.FeedRepository
}

func NewFeedScraperCreator(scraper domain.FeedScraper, extractors []domain.FeedContentExtractor, repository domain.FeedRepository) *FeedScraperCreator {
	return &FeedScraperCreator{
		scraper,
		extractors,
		repository,
	}
}

func (fsc *FeedScraperCreator) Execute() error {
	ctx := context.Background()

	feeds, err := fsc.scraper.Execute(fsc.extractors)
	if err != nil {
		return nil
	}

	for _, feed := range feeds {
		err = fsc.repository.Save(ctx, feed)
		if err != nil {
			return err
		}
	}

	return nil
}
