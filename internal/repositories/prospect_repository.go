package repositories

import (
	"context"
	"fverify_be/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProspectRepositoryImpl struct {
	collection *mongo.Collection
}

func NewProspectRepository(client *mongo.Client, dbName, collectionName string) *ProspectRepositoryImpl {
	collection := client.Database(dbName).Collection(collectionName)
	return &ProspectRepositoryImpl{collection: collection}
}

func (r *ProspectRepositoryImpl) Create(ctx context.Context, prospect *models.Prospect) error {
	_, err := r.collection.InsertOne(ctx, prospect)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProspectRepositoryImpl) GetByID(ctx context.Context, id string) (*models.Prospect, error) {
	var prospect models.Prospect
	err := r.collection.FindOne(ctx, bson.M{"uid": id}).Decode(&prospect)
	return &prospect, err
}

func (r *ProspectRepositoryImpl) Update(ctx context.Context, prospect *models.Prospect) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"uid": prospect.UId}, bson.M{"$set": prospect})
	return err
}

func (r *ProspectRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"uid": id})
	return err
}

func (r *ProspectRepositoryImpl) FindAll(ctx context.Context) ([]*models.Prospect, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prospects []*models.Prospect
	for cursor.Next(ctx) {
		var prospect models.Prospect
		if err := cursor.Decode(&prospect); err != nil {
			return nil, err
		}
		prospects = append(prospects, &prospect)
	}
	return prospects, nil
}
func (r *ProspectRepositoryImpl) GetProspects(ctx context.Context, skip int, limit int) ([]models.Prospect, error) {
	var prospects []models.Prospect

	// MongoDB query with skip and limit
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the results
	if err := cursor.All(ctx, &prospects); err != nil {
		return nil, err
	}

	return prospects, nil
}

func (r *ProspectRepositoryImpl) GetProspectsCount(ctx context.Context) (int, error) {
	// MongoDB query to count documents
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
