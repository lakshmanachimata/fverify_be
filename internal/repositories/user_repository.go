package repositories

import (
	"context"
	"kowtha_be/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName, collectionName string) *UserRepositoryImpl {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserRepositoryImpl{collection: collection}
}

func (r *UserRepositoryImpl) getNextUserID(ctx context.Context) (int, error) {
	counter := struct {
		SequenceValue int `bson:"sequence_value"`
	}{}

	filter := bson.M{"_id": "user_uid"}                                 // Counter identifier for uId
	update := bson.M{"$inc": bson.M{"sequence_value": 1}}               // Increment the counter
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After) // Return the updated document

	// Increment the counter and retrieve the updated value
	err := r.collection.Database().Collection("counters").FindOneAndUpdate(ctx, filter, update, opts).Decode(&counter)
	if err != nil {
		return 0, err
	}

	return counter.SequenceValue, nil
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.UserModel) (*models.UserModel, error) {
	// Generate the next unique ID (uId)
	nextUId, err := r.getNextUserID(ctx)
	if err != nil {
		return nil, err
	}
	user.UId = nextUId // Set the auto-incremented uId

	// Set CreatedTime and UpdatedTime
	user.CreatedTime = time.Now()
	user.UpdatedTime = time.Now()

	// Insert the user into the collection
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

func (r *UserRepositoryImpl) GetByUserID(ctx context.Context, userId string) (*models.UserModel, error) {
	var user models.UserModel
	err := r.collection.FindOne(ctx, bson.M{"userid": userId}).Decode(&user)
	return &user, err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *models.UserModel) error {
	user.UpdatedTime = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"userid": user.UserId}, bson.M{"$set": user})
	return err
}

func (r *UserRepositoryImpl) DeleteByUserId(ctx context.Context, userId string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"userid": userId})
	return err
}

func (r *UserRepositoryImpl) GetAllUsers(ctx context.Context) ([]*models.UserModel, error) {
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
