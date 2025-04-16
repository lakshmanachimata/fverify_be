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
// @Description Organisation model containing all organisation-related information.
//
//	@Example {
//	  "org_id": "12345",
//	  "org_name": "Acme Corp",
//	  "status": "Active"
//	}
type Organisation struct {
	OrgId   string             `json:"org_id" bson:"org_id" example:"12345"`         // Organisation ID
	OrgName string             `json:"org_name" bson:"org_name" example:"Acme Corp"` // Organisation Name
	OrgUUID string             `json:"org_uuid" bson:"org_uuid" example:"uuid-v4"`   // Auto-generated UUID
	Status  OrganisationStatus `json:"status" bson:"status" example:"Active"`        // Organisation Status
}
