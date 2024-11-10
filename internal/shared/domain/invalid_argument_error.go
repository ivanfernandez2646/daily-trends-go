package domain

import "fmt"

type InvalidArgumentError struct {
	message string
}

func NewInvalidArgumentError(message string) InvalidArgumentError {
	return InvalidArgumentError{message}
}

func (iae InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument error: %s", iae.message)
}
