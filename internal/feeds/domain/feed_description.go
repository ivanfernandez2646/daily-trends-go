package domain

type FeedDescription *string

func NewFeedDescription(value string) FeedDescription {
	if value == "" {
		return nil
	}

	return FeedDescription(&value)
}
