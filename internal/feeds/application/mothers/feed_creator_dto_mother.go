package mothers

import (
	"daily-trends/go/internal/feeds/application"
	"daily-trends/go/internal/feeds/domain"
	"math/rand"

	"github.com/bxcodec/faker/v3"
)

func NewFeedCreatorDTO() application.FeedCreatorDTO {
	return application.FeedCreatorDTO{
		Title:       faker.Word(),
		Author:      faker.Name(),
		Description: faker.Paragraph(),
		Source:      domain.FeedSource(rand.Intn(3)).String(),
	}
}

func NewFeedCreatorInvalidDTO() application.FeedCreatorDTO {
	return application.FeedCreatorDTO{}
}
