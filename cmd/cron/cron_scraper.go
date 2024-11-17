package xcron

import (
	"daily-trends/go/internal/feeds/application"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type CronScraper struct {
	scraperCreator application.FeedScraperCreator
}

func NewCronScraper(scraperCreator *application.FeedScraperCreator) *CronScraper {
	return &CronScraper{*scraperCreator}
}

func (cs CronScraper) StartJob() {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatalf("failed to load Madrid timezone: %v", err)
	}

	c := cron.New(cron.WithLocation(location))

	_, err = c.AddFunc("32 17 * * *", func() {
		log.Println("scraper creator job is running...")
		err := cs.scraperCreator.Execute()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("scraper creator job completed successfully")
	})
	if err != nil {
		log.Fatalf("error adding cron job: %v", err)
	}

	c.Start()
}
