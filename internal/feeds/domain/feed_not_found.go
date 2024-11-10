package domain

import "fmt"

type FeedNotFoundError struct {
	id FeedId
}

func NewFeedNotFoundError(id FeedId) FeedNotFoundError {
	return FeedNotFoundError{id}
}

func (fnfe FeedNotFoundError) Error() string {
	return fmt.Sprintf("feed with id %s not found", fnfe.id.Value())
}
