package repositories

import (
	"context"
	"kowtha_be/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName, collectionName string) *UserRepositoryImpl {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserRepositoryImpl{collection: collection}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.UserModel) (*models.UserModel, error) {
	user.CreatedTime = time.Now()
	user.UpdatedTime = time.Now()
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	// Retrieve the inserted user document
	var createdUser models.UserModel
	err = r.collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdUser)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int) (*models.UserModel, error) {
	var user models.UserModel
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	return &user, err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *models.UserModel) error {
	user.UpdatedTime = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"id": user.ID}, bson.M{"$set": user})
	return err
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context) ([]*models.UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.UserModel
	for cursor.Next(ctx) {
		var user models.UserModel
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
