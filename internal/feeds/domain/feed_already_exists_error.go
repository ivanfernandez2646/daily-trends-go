package domain

import "fmt"

type FeedAlreadyExistsError struct {
	id FeedId
}

func NewFeedAlreadyExistsError(id FeedId) FeedAlreadyExistsError {
	return FeedAlreadyExistsError{id}
}

func (faee FeedAlreadyExistsError) Error() string {
	return fmt.Sprintf("feed with id %s already exists", faee.id.Value())
}
