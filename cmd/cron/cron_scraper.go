package xcron

import (
	"daily-trends/go/internal/feeds/application"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

type CronScraper struct {
	scraperCreator application.FeedScraperCreator
}

func NewCronScraper(scraperCreator *application.FeedScraperCreator) *CronScraper {
	return &CronScraper{*scraperCreator}
}

func (cs CronScraper) StartJob() {
	c := cron.New()

	_, err := c.AddFunc("@every 5m", func() {
		err := cs.scraperCreator.Execute()
		if err != nil {
			fmt.Println(err)
		}
	})
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	c.Start()
}
