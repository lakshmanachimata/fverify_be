package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
}

type MongoRepository struct {
	client *mongo.Client
}

func (r *MongoRepository) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	r.client = client
	return nil
}

func (r *MongoRepository) Disconnect(ctx context.Context) error {
	if r.client != nil {
		return r.client.Disconnect(ctx)
	}
	return nil
}
