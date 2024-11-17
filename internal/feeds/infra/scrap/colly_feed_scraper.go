package scrap

import (
	"daily-trends/go/internal/feeds/domain"
	shared_domain "daily-trends/go/internal/shared/domain"
	"sync"

	"fmt"

	"github.com/gocolly/colly/v2"
)

type CollyFeedScraper struct {
	clock shared_domain.Clock
}

func NewCollyFeedScraper(clock shared_domain.Clock) *CollyFeedScraper {
	return &CollyFeedScraper{clock}
}

func (p *CollyFeedScraper) Execute(extractors []domain.FeedContentExtractor) ([]*domain.Feed, error) {
	var res []*domain.Feed
	var sw sync.RWMutex
	var wg sync.WaitGroup

	wg.Add(len(extractors))

	for _, extractor := range extractors {
		go p.processUrl(extractor, &res, &sw, &wg)
	}

	wg.Wait()

	return res, nil
}

func (p *CollyFeedScraper) processUrl(extractor domain.FeedContentExtractor, res *[]*domain.Feed, sw *sync.RWMutex, wg *sync.WaitGroup) {
	defer wg.Done()

	var tmpRes []*domain.Feed

	url := extractor.GetURL()
	selectors := extractor.GetSelectors()
	limit := extractor.GetLimit()

	c := colly.NewCollector()

	c.OnHTML("article", func(e *colly.HTMLElement) {
		if len(tmpRes) >= limit {
			return
		}

		txtAuthor := e.ChildText(selectors.AuthorSelector)
		txtTitle := e.ChildText(selectors.TitleSelector)
		txtDescription := e.ChildText(selectors.DescriptionSelector)

		if txtAuthor == "" || txtTitle == "" || txtDescription == "" {
			return
		}

		uuid := shared_domain.NewRandomUUID()

		data, err := domain.NewFeed(p.clock, domain.WithId(uuid.Value()), domain.WithTitle(txtTitle), domain.WithAuthor(txtAuthor), domain.WithDescription(txtDescription), domain.WithSource(extractor.GetSource().String()))
		if err != nil {
			fmt.Printf("error creating new feed from colly scraper: %v\n", err)
			return
		}

		tmpRes = append(tmpRes, data)
	})

	if selectors.DecodeChars {
		c.OnResponse(func(r *colly.Response) {
			r.Headers.Set("Content-Type", "text/html; charset=iso-8859-1")
		})
	}

	c.Visit(url)
	sw.Lock()
	*res = append(*res, tmpRes...)
	sw.Unlock()
}
