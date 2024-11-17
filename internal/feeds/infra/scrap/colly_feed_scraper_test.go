package scrap_test

import (
	"daily-trends/go/internal/feeds/domain"
	"daily-trends/go/internal/feeds/infra/scrap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollyScraper(t *testing.T) {
	colly := scrap.NewCollyFeedScraper()
	extractors := []domain.FeedContentExtractor{scrap.ElMundoContentExtractor{}, scrap.ElPaisContentExtractor{}}

	var limit int
	for _, extractor := range extractors {
		limit += extractor.GetLimit()
	}

	res, err := colly.Execute(extractors)

	assert.Nil(t, err)
	assert.Equal(t, limit, len(res))
}
