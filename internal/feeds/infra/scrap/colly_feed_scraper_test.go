package scrap_test

import (
	"daily-trends/go/internal/feeds/domain"
	"daily-trends/go/internal/feeds/infra/scrap"
	shared_mocks "daily-trends/go/internal/shared/domain/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollyScraper(t *testing.T) {
	clock := shared_mocks.NewClockMock()
	colly := scrap.NewCollyFeedScraper(clock)
	extractors := []domain.FeedContentExtractor{scrap.ElMundoContentExtractor{}, scrap.ElPaisContentExtractor{}}

	var limit int
	for _, extractor := range extractors {
		limit += extractor.GetLimit()
	}

	res, err := colly.Execute(extractors)

	assert.Nil(t, err)
	assert.Equal(t, limit, len(res))
}
