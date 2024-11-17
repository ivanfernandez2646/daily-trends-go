package scrap

import (
	"crypto/tls"
	"daily-trends/go/internal/feeds/domain"
	shared_domain "daily-trends/go/internal/shared/domain"
	"log"
	"net/http"
	"sync"

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

	// Disable SSL (only test and personal use)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnHTML("article", func(e *colly.HTMLElement) {
		if len(tmpRes) >= limit {
			return
		}

		txtAuthor := e.ChildText(selectors.AuthorSelector)
		txtTitle := e.ChildText(selectors.TitleSelector)
		txtDescription := e.ChildText(selectors.DescriptionSelector)
		txtUrl := e.ChildAttr("a", "href")

		if txtAuthor == "" || txtTitle == "" || txtDescription == "" || txtUrl == "" {
			return
		}

		uuid := shared_domain.NewRandomUUID()

		data, err := domain.NewFeed(p.clock, domain.WithId(uuid.Value()), domain.WithTitle(txtTitle), domain.WithAuthor(txtAuthor), domain.WithDescription(txtDescription), domain.WithSource(extractor.GetSource().String()), domain.WithUrl(txtUrl))
		if err != nil {
			log.Printf("error creating new feed from colly scraper: %v\n", err)
			return
		}

		tmpRes = append(tmpRes, data)
	})

	if selectors.DecodeChars {
		c.OnResponse(func(r *colly.Response) {
			r.Headers.Set("Content-Type", "text/html; charset=iso-8859-1")
		})
	}

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	sw.Lock()
	*res = append(*res, tmpRes...)
	sw.Unlock()
}
