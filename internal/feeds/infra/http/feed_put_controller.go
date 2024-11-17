package xhttp

import (
	"daily-trends/go/internal/feeds/application"
	"daily-trends/go/internal/feeds/domain"
	shared_domain "daily-trends/go/internal/shared/domain"
	"log"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewPutFeedController(creator *application.FeedCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var createFeedDto application.FeedCreatorDTO
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		if err := decoder.Decode(&createFeedDto); err != nil {
			http.Error(w, fmt.Sprintf("error decoding body: %v", err), http.StatusBadRequest)
			return
		}

		err := creator.Execute(id, &createFeedDto)
		if err != nil {
			log.Println(err)

			w.Header().Set("Content-Type", "application/json")

			switch err.(type) {
			case shared_domain.InvalidArgumentError:
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error": "%s"}`, err)
				return
			case domain.FeedAlreadyExistsError:
				w.WriteHeader(http.StatusConflict)
				fmt.Fprintf(w, `{"error": "%s"}`, err)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error": "internal server error"}`)
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
	}
}
