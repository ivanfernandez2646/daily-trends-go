package domain

import (
	shared_domain "daily-trends/go/internal/shared/domain"
	"fmt"
)

type FeedSource int

const (
	CMS = iota
	EL_PAIS
	EL_MUNDO
)

func (s FeedSource) String() string {
	switch s {
	case CMS:
		return "CMS"
	case EL_PAIS:
		return "EL_PAIS"
	case EL_MUNDO:
		return "EL_MUNDO"
	default:
		return "unknown source"
	}
}

func NewFeedSource(value string) (FeedSource, error) {
	switch value {
	case "CMS":
		return CMS, nil
	case "EL_PAIS":
		return EL_PAIS, nil
	case "EL_MUNDO":
		return EL_MUNDO, nil
	default:
		return -1, shared_domain.NewInvalidArgumentError(fmt.Sprintf("feed source invalid value %s", value))
	}
}
