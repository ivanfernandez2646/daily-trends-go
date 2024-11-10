package domain

type FeedContentSelectors struct {
	AuthorSelector      string
	TitleSelector       string
	DescriptionSelector string
	DecodeChars         bool
}

type FeedContentExtractor interface {
	GetSelectors() FeedContentSelectors
	GetURL() string
	GetSource() FeedSource
}

type FeedScraper interface {
	Execute(extractors []FeedContentExtractor) ([]*Feed, error)
}
