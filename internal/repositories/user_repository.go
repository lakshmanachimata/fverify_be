package repositories

import (
	"context"
	"fverify_be/internal/models"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}

// HashPassword hashes a plain-text password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with a plain-text password
func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func NewUserRepository(client *mongo.Client, dbName, collectionName string) *UserRepositoryImpl {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserRepositoryImpl{collection: collection}
}

func (r *UserRepositoryImpl) ValidateUser(ctx context.Context, username, password string, orgUUID string) (*models.UserModel, error) {
	var user models.UserModel
	err := r.collection.FindOne(ctx, bson.M{"username": username, "org_uuid": orgUUID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	// Compare the hashed password with the provided password
	if err := CheckPassword(user.Password, password); err != nil {
		return nil, err // Password does not match
	}

	return &user, nil
}

func (r *UserRepositoryImpl) SetPassword(ctx context.Context, uId string, newPassword string) error {
	// Hash the new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update the password for the user with the given uId
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"uid": uId}, // Filter by uId
		bson.M{"$set": bson.M{"password": hashedPassword}}, // Update the password field
	)
	return err
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.UserModel) (*models.UserRespModel, error) {
	// Generate the next unique ID (uId)
	// Hash the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword // Set the hashed password
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

	// Convert to UserRespModel
	createdUserResp := models.UserRespModel{
		UId:           createdUser.UId,
		UserId:        createdUser.UserId,
		Username:      createdUser.Username,
		Role:          createdUser.Role,
		Status:        createdUser.Status,
		MobileNumber:  createdUser.MobileNumber,
		Remarks:       createdUser.Remarks,
		OrgUUID:       createdUser.OrgUUID,
		CreatedTime:   createdUser.CreatedTime,
		OrgStatus:     createdUser.OrgStatus,
		UpdatedTime:   createdUser.UpdatedTime,
		UpdateHistory: createdUser.UpdateHistory}
	return &createdUserResp, nil
}

func (r *UserRepositoryImpl) GetByUserID(ctx context.Context, userId string) (*models.UserRespModel, error) {
	var user models.UserRespModel
	err := r.collection.FindOne(ctx, bson.M{"userid": userId}).Decode(&user)
	return &user, err
}

func (r *UserRepositoryImpl) GetByUserUID(ctx context.Context, uid string) (*models.UserRespModel, error) {
	var user models.UserRespModel
	err := r.collection.FindOne(ctx, bson.M{"uid": uid}).Decode(&user)
	return &user, err
}

func (r *UserRepositoryImpl) DeleteByUId(ctx context.Context, uId string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"uid": uId})
	return err
}
func (r *UserRepositoryImpl) DeleteByUserId(ctx context.Context, userId string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"userid": userId})
	return err
}

func (r *UserRepositoryImpl) GetAllUsers(ctx context.Context) ([]*models.UserRespModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.UserRespModel
	for cursor.Next(ctx) {
		var user models.UserRespModel
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *models.UserModel, authUserName string) (*models.UserRespModel, error) {
	// Update the UpdatedTime field
	var eUser models.UserModel
	err := r.collection.FindOne(ctx, bson.M{"uid": user.UId}).Decode(&eUser)
	if err != nil {
		return nil, err
	}

	// Generate a diff between the existing user and the incoming user
	var updateComments []string
	if eUser.UserId != user.UserId {
		updateComments = append(updateComments, "user id  changed from '"+eUser.UserId+"' to '"+user.UserId+"'")
	}
	if eUser.Username != user.Username {
		updateComments = append(updateComments, "user name changed from '"+eUser.Username+"' to '"+user.Username+"'")
	}
	if eUser.Password != user.Password {
		updateComments = append(updateComments, "password updated")
	}
	if eUser.Role != user.Role {
		updateComments = append(updateComments, "role changed from '"+string(eUser.Role)+"' to '"+string(user.Role)+"'")
	}
	if eUser.Status != user.Status {
		updateComments = append(updateComments, "status changed from '"+string(eUser.Status)+"' to '"+string(user.Status)+"'")
	}
	if eUser.Remarks != user.Remarks {
		updateComments = append(updateComments, "remarks changed from '"+eUser.Remarks+"' to '"+user.Remarks+"'")
	}
	if eUser.MobileNumber != user.MobileNumber {
		updateComments = append(updateComments, "mobile number changed from '"+eUser.MobileNumber+"' to '"+user.MobileNumber+"'")
	}

	user.UpdatedTime = time.Now().UTC().Format(time.RFC3339)
	user.UpdateHistory = append(user.UpdateHistory, models.UpdateHistory{
		UpdatedTime:     time.Now().UTC().Format(time.RFC3339),
		UpdatedComments: strings.Join(updateComments, ", "),
		UpdateBy:        authUserName,
	})

	// Perform the update operation
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"uid": user.UId}, // Filter by uId
		bson.M{"$set": user},    // Update the user document
	)
	// Convert to UserRespModel
	createdUserResp := models.UserRespModel{
		UId:           user.UId,
		UserId:        user.UserId,
		Username:      user.Username,
		Role:          user.Role,
		Status:        user.Status,
		MobileNumber:  user.MobileNumber,
		Remarks:       user.Remarks,
		OrgUUID:       user.OrgUUID,
		CreatedTime:   user.CreatedTime,
		OrgStatus:     user.OrgStatus,
		UpdatedTime:   user.UpdatedTime,
		UpdateHistory: user.UpdateHistory}
	return &createdUserResp, nil
}
func (r *UserRepositoryImpl) UpdateUsersStatusByOrgUUID(ctx context.Context, orgUUID string, status models.UserStatus) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"org_uuid": orgUUID},              // Filter by org_uuid
		bson.M{"$set": bson.M{"status": status}}, // Update the status field
	)
	return err
}

func (r *UserRepositoryImpl) UpdateUserStatus(ctx context.Context, userId string, status string) error {
	filter := bson.M{"userid": userId}
	update := bson.M{"$set": bson.M{"status": status, "updated_time": time.Now().UTC().Format(time.RFC3339)}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
