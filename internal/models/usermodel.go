package models

import "time"

type Role string
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

type UpdateHistory struct {
	UpdatedComments string    `bson:"updated_comments" json:"updated_comments"`
	UpdatedTime     time.Time `bson:"updated_time" json:"updated_time"`
}

type UserModel struct {
	ID            int             `bson:"user_id" json:"user_id"` // Incremental ID
	Username      string          `bson:"username" json:"username"`
	Password      string          `bson:"password" json:"password"`
	Role          Role            `bson:"role" json:"role"`
	Status        UserStatus      `bson:"status" json:"status"`
	CreatedTime   time.Time       `bson:"created_time" json:"created_time"`
	UpdatedTime   time.Time       `bson:"updated_time" json:"updated_time"`
	UpdateHistory []UpdateHistory `bson:"update_history" json:"update_history"`
	Remarks       string          `bson:"remarks" json:"remarks"`
}
