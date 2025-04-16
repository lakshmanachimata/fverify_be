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
	UpdateBy        string    `bson:"update_by" json:"update_by" example:"admin"`                           // User who made the update
}

// UserModel represents a user in the system.
// @Description User model containing all user-related information.
//
//		@Example {
//		  "userid": "12345",
//		  "username": "john_doe",
//		  "password": "hashed_password",
//		  "role": "Admin",
//		  "status": "Active",
//		  "created_time": "2023-04-12T15:04:05Z",
//		  "updated_time": "2023-04-12T15:04:05Z",
//	   	  "mobile_number": "9876543210",
//		  "update_history": [
//		    {
//		      "updated_comments": "Updated user role",
//		      "updated_time": "2023-04-12T15:04:05Z"
//		    }
//		  ],
//		  "remarks": "User is active and verified",
//		  "orgId": "123456"
//		}
type UserModel struct {
	UId           string             `bson:"uid" json:"uid" example:"123e4567-e89b-12d3-a456-426614174111"`           // Auto-incremented unique identifier
	UserId        string             `bson:"userid" json:"userid" example:"112345"`                                   // Unique identifier for the user
	Username      string             `bson:"username" json:"username" example:"john_doe"`                             // Username of the user
	Password      string             `bson:"password" json:"password" example:"hashed_password"`                      // Hashed password
	Role          Role               `bson:"role" json:"role" example:"Admin"`                                        // Role of the user
	Status        UserStatus         `bson:"status" json:"status" example:"Active"`                                   // Status of the user
	CreatedTime   time.Time          `bson:"created_time" json:"created_time" example:"2023-04-12T15:04:05Z"`         // Time when the user was created
	UpdatedTime   time.Time          `bson:"updated_time" json:"updated_time" example:"2023-04-12T15:04:05Z"`         // Time when the user was last updated
	UpdateHistory []UpdateHistory    `bson:"update_history" json:"update_history"`                                    // History of updates
	Remarks       string             `bson:"remarks" json:"remarks" example:"User is active and verified"`            // Additional remarks about the user
	MobileNumber  string             `bson:"mobile_number" json:"mobile_number" example:"9876543210"`                 // Mobile number of the user
	OrgStatus     OrganisationStatus `bson:"org_status" json:"org_status" example:"123456"`                           // Organization ID
	OrgUUID       string             `bson:"org_uuid" json:"org_uuid" example:"123e4567-e89b-12d3-a456-426614174000"` // UUID of the organization
}

// LoginRequest represents the request payload for the login API.
// @Description Login request payload containing username and password.
//
//	@Example {
//	  "username": "john_doe",
//	  "password": "password"
//	}
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"` // Username
	Password string `json:"password" binding:"required" example:"password"` // Password
	OrgId    string `json:"orgId" binding:"required" example:"123456"`      // Organization ID
}

// LoginResponse represents the response payload for the login API.
// @Description Login response payload containing user details and the generated token.
//
//	@Example {
//	  "uid": 1,
//	  "userId": "12345",
//	  "username": "john_doe",
//	  "role": "Admin",
//	  "status": "Active",
//	  "mobileNumber": "9876543210",
//	  "token": "<jwt_token>"
//	}
type LoginResponse struct {
	UId          string `json:"uid" example:"1"`                   // User's unique ID
	UserId       string `json:"userId" example:"12345"`            // User's unique identifier
	Username     string `json:"username" example:"john_doe"`       // Username
	Role         string `json:"role" example:"Admin"`              // Role of the user
	Status       string `json:"status" example:"Active"`           // Status of the user
	MobileNumber string `json:"mobileNumber" example:"9876543210"` // Mobile number
	Token        string `json:"token" example:"<jwt_token>"`       // Auth token
}

type SetPasswordRequest struct {
	Password string `json:"password" binding:"required" example:"new_password"` // New password
}
