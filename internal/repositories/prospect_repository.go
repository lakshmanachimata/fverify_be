package repositories

import (
	"context"
	"kowtha_be/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProspectRepositoryImpl struct {
	collection *mongo.Collection
}

func NewProspectRepository(client *mongo.Client, dbName, collectionName string) *ProspectRepositoryImpl {
	collection := client.Database(dbName).Collection(collectionName)
	return &ProspectRepositoryImpl{collection: collection}
}

func (r *ProspectRepositoryImpl) Create(ctx context.Context, prospect *models.ProspectModel) error {
	_, err := r.collection.InsertOne(ctx, prospect)
	return err
}

func (r *ProspectRepositoryImpl) GetByID(ctx context.Context, id string) (*models.ProspectModel, error) {
	var prospect models.ProspectModel
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&prospect)
	return &prospect, err
}

func (r *ProspectRepositoryImpl) Update(ctx context.Context, prospect *models.ProspectModel) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"id": prospect.ID}, bson.M{"$set": prospect})
	return err
}

func (r *ProspectRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *ProspectRepositoryImpl) FindAll(ctx context.Context) ([]*models.ProspectModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prospects []*models.ProspectModel
	for cursor.Next(ctx) {
		var prospect models.ProspectModel
		if err := cursor.Decode(&prospect); err != nil {
			return nil, err
		}
		prospects = append(prospects, &prospect)
	}
	return prospects, nil
}
