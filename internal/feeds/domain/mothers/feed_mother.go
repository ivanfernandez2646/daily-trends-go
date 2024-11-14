package mothers

import (
	"daily-trends/go/internal/feeds/application"
	"daily-trends/go/internal/feeds/domain"

	"math/rand"

	"github.com/bxcodec/faker/v3"
)

func NewRandomFeed() (*domain.Feed, error) {
	id := faker.UUIDDigit()
	title := faker.Word()
	author := faker.Name()
	source := domain.FeedSource(rand.Intn(3))

	description := ""
	if rand.Intn(2) == 0 {
		description = faker.Word()
	}

	return domain.NewFeed(domain.WithId(id), domain.WithTitle(title), domain.WithDescription(description), domain.WithAuthor(author), domain.WithSource(source.String()))
}

func NewFeedFromProps(id string, dto application.FeedCreatorDTO) (*domain.Feed, error) {
	return domain.NewFeed(domain.WithId(id), domain.WithTitle(dto.Title), domain.WithDescription(dto.Description), domain.WithAuthor(dto.Author), domain.WithSource(dto.Source))
}
