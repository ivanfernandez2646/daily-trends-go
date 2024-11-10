package xhttp

import (
	"daily-trends/go/internal/feeds/application"
	"fmt"
	"net/http"
)

func NewFeedScraperGetController(scraperCreator *application.FeedScraperCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := scraperCreator.Execute()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "internal server error"}`)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
