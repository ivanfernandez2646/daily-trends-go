package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type UUID struct {
	value string
}

func NewUUID(value string) (*UUID, error) {
	val, err := uuid.Parse(value)
	if err != nil {
		return nil, NewInvalidArgumentError(fmt.Sprintf("cannot parse uuid %s", value))
	}

	return &UUID{
		value: val.String(),
	}, nil
}

func NewRandomUUID() UUID {
	return UUID{
		value: uuid.NewString(),
	}
}

func (u UUID) Value() string {
	return u.value
}
