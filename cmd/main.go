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
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	// DI
	repo, err := persistence.NewMongoDBFeedRepository(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	clock := infra.NewClockworkClock()
	scraper := scrap.NewCollyFeedScraper(clock)
	creator := application.NewFeedCreator(repo, clock)
	finder := application.NewFeedFinder(repo)
	searcher := application.NewFeedSearcher(repo)
	scraperCreator := application.NewFeedScraperCreator(scraper, []domain.FeedContentExtractor{scrap.ElPaisContentExtractor{}, scrap.ElMundoContentExtractor{}}, repo)

	// Starts crons
	cs := xcron.NewCronScraper(scraperCreator)
	cs.StartJob()

	// HTTP
	r := mux.NewRouter()

	// Add cors
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	r.Use(corsHandler)

	// Routes
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Println("health ok")
		fmt.Fprint(w, `{status:"ok"}`)
	}).Methods("GET")
	r.HandleFunc("/feeds/home", xhttp.NewGetFeedHomeController(searcher)).Methods("GET")
	r.HandleFunc("/feeds/{id}", xhttp.NewPutFeedController(creator)).Methods("PUT")
	r.HandleFunc("/feeds/{id}", xhttp.NewGetFeedController(finder)).Methods("GET")

	// Initialize server
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("error parsing PORT: ", err)
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", host, strconv.Itoa(port)),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start server
	log.Printf("starting server on %s ...\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func loadEnv() {
	appEnv := os.Getenv("APP_ENV")

	if appEnv != "development" {
		return
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
}
