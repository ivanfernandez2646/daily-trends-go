package main

import (
	"context"
	xcron "daily-trends/go/cmd/cron"
	"daily-trends/go/internal/feeds/application"
	"daily-trends/go/internal/feeds/domain"
	xhttp "daily-trends/go/internal/feeds/infra/http"
	"daily-trends/go/internal/feeds/infra/persistence"
	"daily-trends/go/internal/feeds/infra/scrap"
	"daily-trends/go/internal/shared/infra"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	HOST = "localhost"
	PORT = 9292
)

func main() {
	// DI
	repo, err := persistence.NewMongoDBFeedRepository(context.Background())
	if err != nil {
		panic(err)
	}
	clock := infra.NewClockworkClock()
	scraper := scrap.NewCollyFeedScraper(clock)
	creator := application.NewFeedCreator(repo, clock)
	finder := application.NewFeedFinder(repo)
	scraperCreator := application.NewFeedScraperCreator(scraper, []domain.FeedContentExtractor{scrap.ElPaisContentExtractor{}, scrap.ElMundoContentExtractor{}}, repo)

	// Start cron scraper job
	cs := xcron.NewCronScraper(scraperCreator)
	cs.StartJob()

	// Register HTTP routes
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{status:"ok"}`)
	}).Methods("GET")
	r.HandleFunc("/feeds/{id}", xhttp.NewPutFeedController(creator)).Methods("PUT")
	r.HandleFunc("/feeds/{id}", xhttp.NewGetFeedController(finder)).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", HOST, strconv.Itoa(PORT)),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start server
	fmt.Printf("starting server on %s ...\n", srv.Addr)
	srv.ListenAndServe()
}
