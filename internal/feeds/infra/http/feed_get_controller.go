package xhttp

import (
	"daily-trends/go/internal/feeds/application"
	"daily-trends/go/internal/feeds/domain"
	"daily-trends/go/internal/feeds/infra/http/responses"
	shared_domain "daily-trends/go/internal/shared/domain"
	"log"

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
			log.Println(err)

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

		resDto := responses.NewFeedGetResponse(feed)
		bytes, err := json.Marshal(resDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%s"}`, err)
		}
		w.Write(bytes)
	}
}
