package xhttp

import (
	"daily-trends/go/internal/feeds/application"
	"daily-trends/go/internal/feeds/infra/http/responses"
	shared_domain "daily-trends/go/internal/shared/domain"
	"log"

	"encoding/json"
	"fmt"
	"net/http"
)

func NewGetFeedHomeController(searcher *application.FeedSearcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		feeds, err := searcher.Execute()
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			log.Println(err)

			switch err.(type) {
			case shared_domain.InvalidArgumentError:
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error": "%s"}`, err)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error": "internal server error"}`)
				return
			}
		}

		resDto := responses.NewFeedHomeGetResponse(feeds)
		bytes, err := json.Marshal(resDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%s"}`, err)
		}
		w.Write(bytes)
	}
}
