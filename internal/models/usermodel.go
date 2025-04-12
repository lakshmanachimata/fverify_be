package models

import "time"

// Role represents the role of a user.
// Enum: "Admin", "Operations Lead", "Field Lead", "Field Executive", "Owner", "Operations Executive"
type Role string

// UserStatus represents the status of a user.
// Enum: "Created", "Confirmed", "Verified", "Active", "Inactive", "Disabled", "Banned"
type UserStatus string

const (
	// Role Enums
	Admin               Role = "Admin"
	OperationsLead      Role = "Operations Lead"
	FieldLead           Role = "Field Lead"
	FieldExecutive      Role = "Field Executive"
	Owner               Role = "Owner"
	OperationsExecutive Role = "Operations Executive"

	// Status Enums
	Created   UserStatus = "Created"
	Confirmed UserStatus = "Confirmed"
	Verified  UserStatus = "Verified"
	Active    UserStatus = "Active"
	InActive  UserStatus = "Inactive"
	Disabled  UserStatus = "Disabled"
	Banned    UserStatus = "Banned"
)

// UpdateHistory represents the history of updates made to a user.
// @Description History of updates made to a user.
type UpdateHistory struct {
	UpdatedComments string    `bson:"updated_comments" json:"updated_comments" example:"Updated user role"` // Comments about the update
	UpdatedTime     time.Time `bson:"updated_time" json:"updated_time" example:"2023-04-12T15:04:05Z"`      // Time of the update
}

// UserModel represents a user in the system.
// @Description User model containing all user-related information.
//
//	@Example {
//	  "user_id": 1,
//	  "username": "john_doe",
//	  "password": "hashed_password",
//	  "role": "Admin",
//	  "status": "Active",
//	  "created_time": "2023-04-12T15:04:05Z",
//	  "updated_time": "2023-04-12T15:04:05Z",
//	  "update_history": [
//	    {
//	      "updated_comments": "Updated user role",
//	      "updated_time": "2023-04-12T15:04:05Z"
//	    }
//	  ],
//	  "remarks": "User is active and verified"
//	}
type UserModel struct {
	ID            int             `bson:"user_id" json:"user_id" example:"1"`                              // Incremental ID
	Username      string          `bson:"username" json:"username" example:"john_doe"`                     // Username of the user
	Password      string          `bson:"password" json:"password" example:"hashed_password"`              // Hashed password
	Role          Role            `bson:"role" json:"role" example:"Admin"`                                // Role of the user
	Status        UserStatus      `bson:"status" json:"status" example:"Active"`                           // Status of the user
	CreatedTime   time.Time       `bson:"created_time" json:"created_time" example:"2023-04-12T15:04:05Z"` // Time when the user was created
	UpdatedTime   time.Time       `bson:"updated_time" json:"updated_time" example:"2023-04-12T15:04:05Z"` // Time when the user was last updated
	UpdateHistory []UpdateHistory `bson:"update_history" json:"update_history"`                            // History of updates
	Remarks       string          `bson:"remarks" json:"remarks" example:"User is active and verified"`    // Additional remarks about the user
}
