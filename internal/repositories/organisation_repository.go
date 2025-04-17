package repositories

import (
	"context"
	"fverify_be/internal/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrganisationRepository struct {
	collection *mongo.Collection
}

func NewOrganisationRepository(client *mongo.Client, dbName, collectionName string) *OrganisationRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &OrganisationRepository{collection: collection}
}

func (r *OrganisationRepository) Create(ctx context.Context, org *models.Organisation) (*models.Organisation, error) {
	// Generate a UUID for the organisation
	org.OrgUUID = uuid.New().String()

	_, err := r.collection.InsertOne(ctx, org)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (r *OrganisationRepository) Update(ctx context.Context, org_id string, org *models.Organisation) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"org_id": org_id},
		bson.M{"$set": org},
	)
	return err
}

func (r *OrganisationRepository) Delete(ctx context.Context, org_id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"org_id": org_id})
	return err
}
func (r *OrganisationRepository) GetAllOrganisations(ctx context.Context) ([]*models.Organisation, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var organisations []*models.Organisation
	for cursor.Next(ctx) {
		var org models.Organisation
		if err := cursor.Decode(&org); err != nil {
			return nil, err
		}
		organisations = append(organisations, &org)
	}
	return organisations, nil
}
func (r *OrganisationRepository) IsOrgActive(ctx context.Context, org_id string) (bool, error) {
	var org models.Organisation
	err := r.collection.FindOne(ctx, bson.M{"org_id": org_id, "status": models.Active}).Decode(&org)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (r *OrganisationRepository) GetOrganisationByID(ctx context.Context, org_id string) (*models.Organisation, error) {
	var org models.Organisation
	err := r.collection.FindOne(ctx, bson.M{"org_id": org_id}).Decode(&org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}
