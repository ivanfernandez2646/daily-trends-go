package xhttp

import (
	"daily-trends/go/internal/feeds/application"
	"daily-trends/go/internal/feeds/domain"
	shared_domain "daily-trends/go/internal/shared/domain"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewGetFeedController(finder *application.FeedFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		feed, err := finder.Execute(id)
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			fmt.Println(err)

			switch err.(type) {
			case shared_domain.InvalidArgumentError:
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error": "%s"}`, err)
				return
			case domain.FeedNotFoundError:
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, `{"error": "%s"}`, err)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error": "internal server error"}`)
				return
			}
		}

		resDto := newFeedGetResponse(feed)
		bytes, err := json.Marshal(resDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%s"}`, err)
		}
		w.Write(bytes)
	}
}

type getFeedResponseDTO struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Author      string  `json:"author"`
	Source      string  `json:"source"`
}

func newFeedGetResponse(feed *domain.Feed) *getFeedResponseDTO {
	id := feed.Id()
	feedDTO := &getFeedResponseDTO{
		Id:          id.Value(),
		Title:       string(feed.Title()),
		Description: feed.Description(),
		Author:      string(feed.Author()),
		Source:      feed.Source().String(),
	}

	return feedDTO
}
