package persistence

import (
	"context"
	"daily-trends/go/internal/feeds/domain"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDBFeedRepository struct {
	client *mongo.Client
}

func NewMongoDBFeedRepository(ctx context.Context) (*MongoDBFeedRepository, error) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://root:secret@127.0.0.1:27017/?authSource=admin"))

	if err != nil {
		return nil, err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err = client.Ping(ctxTimeout, readpref.Primary())
	if err != nil {
		client.Disconnect(ctx)
		return nil, err
	}

	return &MongoDBFeedRepository{client}, nil
}

func (r *MongoDBFeedRepository) Save(ctx context.Context, feed *domain.Feed) error {
	mCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := r.client.Database("daily_trends").Collection("feeds")
	feedBson, err := r.toBSON(feed)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(mCtx, feedBson)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoDBFeedRepository) FindById(ctx context.Context, id domain.FeedId) (*domain.Feed, error) {
	mCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result bson.M
	collection := r.client.Database("daily_trends").Collection("feeds")
	err := collection.FindOne(mCtx, bson.M{"_id": id.Value()}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	feed, err := r.fromBSON(result)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (r *MongoDBFeedRepository) toBSON(feed *domain.Feed) (bson.M, error) {
	id := feed.Id()
	feedBSON := bson.M{
		"_id":         id.Value(),
		"title":       feed.Title(),
		"description": feed.Description(),
		"author":      feed.Author(),
		"source":      feed.Source().String(),
		"createdAt":   feed.CreatedAt().Format(time.RFC3339),
	}

	if feed.Description() == nil {
		delete(feedBSON, "description")
	}

	return feedBSON, nil
}

func (r *MongoDBFeedRepository) fromBSON(value bson.M) (*domain.Feed, error) {
	id, ok := value["_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid _id type")
	}

	title, ok := value["title"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid title type")
	}

	var description string
	if val, exists := value["description"]; exists {
		if str, ok := val.(string); ok {
			description = str
		}
	}

	author, ok := value["author"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid author type")
	}

	source, ok := value["source"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid source type")
	}

	createdAt, ok := value["createdAt"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid createdAt type")
	}

	feed, err := domain.NewFeedOnlyFromOptions(
		domain.WithId(id),
		domain.WithTitle(title),
		domain.WithDescription(description),
		domain.WithAuthor(author),
		domain.WithSource(source),
		domain.WithCreatedAt(createdAt),
	)

	if err != nil {
		return nil, fmt.Errorf("error creating feed from BSON: %v", err)
	}

	return feed, nil
}
