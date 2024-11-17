package scrap

import "daily-trends/go/internal/feeds/domain"

type ElPaisContentExtractor struct{}

func (s ElPaisContentExtractor) GetSelectors() domain.FeedContentSelectors {
	return domain.FeedContentSelectors{
		AuthorSelector:      ".c_a",
		TitleSelector:       ".c_h",
		DescriptionSelector: ".c_d",
	}
}

func (s ElPaisContentExtractor) GetURL() string {
	return "https://elpais.com"
}

func (s ElPaisContentExtractor) GetSource() domain.FeedSource {
	return domain.FeedSource(domain.EL_PAIS)
}

func (s ElPaisContentExtractor) GetLimit() int {
	return 20
}
