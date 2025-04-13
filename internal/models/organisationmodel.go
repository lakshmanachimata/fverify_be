package models

// OrganisationStatus represents the status of an organisation.
// Enum: "Created", "Active", "InActive"
type OrganisationStatus string

const (
	OrgCreated  OrganisationStatus = "Created"
	OrgActive   OrganisationStatus = "Active"
	OrgInActive OrganisationStatus = "InActive"
)

// Organisation represents an organisation in the system.
type Organisation struct {
	OrgId   string             `json:"orgId" bson:"orgId" example:"12345"`         // Organisation ID
	OrgName string             `json:"orgName" bson:"orgName" example:"Acme Corp"` // Organisation Name
	OrgUUID string             `json:"orgUUID" bson:"orgUUID" example:"uuid-v4"`   // Auto-generated UUID
	Status  OrganisationStatus `json:"status" bson:"status" example:"Active"`      // Organisation Status
}
