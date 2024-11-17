package scrap

import "daily-trends/go/internal/feeds/domain"

type ElMundoContentExtractor struct{}

func (s ElMundoContentExtractor) GetSelectors() domain.FeedContentSelectors {
	return domain.FeedContentSelectors{
		AuthorSelector:      ".ue-c-cover-content__byline-name .ue-c-cover-content__link",
		TitleSelector:       ".ue-c-cover-content__headline",
		DescriptionSelector: ".ue-c-cover-content__kicker",
		DecodeChars:         true,
	}
}

func (s ElMundoContentExtractor) GetURL() string {
	return "https://elmundo.es"
}

func (s ElMundoContentExtractor) GetSource() domain.FeedSource {
	return domain.FeedSource(domain.EL_MUNDO)
}

func (s ElMundoContentExtractor) GetLimit() int {
	return 5
}
