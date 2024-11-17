package persistence

import (
	"context"
	"crypto/tls"
	"daily-trends/go/internal/feeds/domain"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDBFeedRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoDBFeedRepository(ctx context.Context) (*MongoDBFeedRepository, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		panic("MONGO_URI invalid value")
	}

	database := os.Getenv("MONDO_DATABASE")
	if database == "" {
		database = "daily_trends"
	}

	clientOptions := options.Client().ApplyURI(uri)

	if os.Getenv("APP_ENV") != "development" {
		clientOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil, err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = client.Ping(ctxTimeout, readpref.Primary())
	if err != nil {
		client.Disconnect(ctx)
		return nil, err
	}

	return &MongoDBFeedRepository{client: client, database: database, collection: "feeds"}, nil
}

func (r *MongoDBFeedRepository) Save(ctx context.Context, feed *domain.Feed) error {
	mCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := r.client.Database(r.database).Collection(r.collection)
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
	collection := r.client.Database(r.database).Collection(r.collection)
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

func (r *MongoDBFeedRepository) Search(ctx context.Context) ([]*domain.Feed, error) {
	mCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var feeds []*domain.Feed

	collection := r.client.Database(r.database).Collection(r.collection)

	findOptions := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(10)
	cursor, err := collection.Find(mCtx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(mCtx)

	for cursor.Next(mCtx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		feed, err := r.fromBSON(result)
		if err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (r *MongoDBFeedRepository) toBSON(feed *domain.Feed) (bson.M, error) {
	id := feed.Id()
	feedBSON := bson.M{
		"_id":         id.Value(),
		"title":       feed.Title(),
		"description": feed.Description(),
		"author":      feed.Author(),
		"source":      feed.Source().String(),
		"url":         feed.Url(),
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

	var url string
	if val, exists := value["url"]; exists {
		if str, ok := val.(string); ok {
			url = str
		}
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
		domain.WithUrl(url),
		domain.WithCreatedAt(createdAt),
	)

	if err != nil {
		return nil, fmt.Errorf("error creating feed from BSON: %v", err)
	}

	return feed, nil
}
