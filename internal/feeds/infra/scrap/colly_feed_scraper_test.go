package scrap_test

import (
	"daily-trends/go/internal/feeds/domain"
	"daily-trends/go/internal/feeds/infra/scrap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollyScraper(t *testing.T) {
	colly := scrap.NewCollyFeedScraper()

	res, err := colly.Execute([]domain.FeedContentExtractor{scrap.ElMundoContentExtractor{}, scrap.ElPaisContentExtractor{}})

	assert.Nil(t, err)
	assert.Equal(t, 40, len(res))
}
