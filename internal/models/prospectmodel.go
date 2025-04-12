package models

// ProspectStatus represents the status of a prospect.
// Enum: "Pending", "OnVisit", "Progress", "Approved", "Rejected", "UnderReview", "Completed", "Submitted", "Cancelled", "RePending", "Postponed"
type ProspectStatus string

const (
	Pending     ProspectStatus = "Pending"
	OnVisit     ProspectStatus = "OnVisit"
	Progressve  ProspectStatus = "Progress"
	Approved    ProspectStatus = "Approved"
	Rejected    ProspectStatus = "Rejected"
	UnderReview ProspectStatus = "UnderReview"
	Completed   ProspectStatus = "Completed"
	Submitted   ProspectStatus = "Submitted"
	Cancelled   ProspectStatus = "Cancelled"
	RePending   ProspectStatus = "RePending"
	Postponed   ProspectStatus = "Postponed"
)

// ProspectModel represents a prospect in the system.
// @Description Prospect model containing all prospect-related information.
//
//	@Example {
//	  "id": 1,
//	  "prospect_id": "P12345",
//	  "applicant_name": "John Doe",
//	  "mobile_number": "9876543210",
//	  "gender": "Male",
//	  "age": 30,
//	  "residential_address": "123 Main Street",
//	  "years_of_stay": 5,
//	  "number_of_family_members": 4,
//	  "reference_name": "Jane Doe",
//	  "reference_relation": "Sister",
//	  "reference_mobile": "9876543211",
//	  "employment_type": "Employee",
//	  "office_address": "456 Office Street",
//	  "years_in_current_office": 3,
//	  "role": "Manager",
//	  "emp_id": "EMP123",
//	  "status": "Pending",
//	  "previous_experience": "5 years in sales",
//	  "gross_salary": 50000.00,
//	  "net_salary": 40000.00,
//	  "colleague_name": "Mark Smith",
//	  "colleague_designation": "Team Lead",
//	  "colleague_mobile": "9876543212",
//	  "uploaded_images": ["image1.jpg", "image2.jpg"],
//	  "remarks": "Prospect is under review"
//	}
type ProspectModel struct {
	ID                    int            `bson:"id" json:"id" example:"1"`                                                          // Incremental ID
	ProspectId            string         `bson:"prospect_id" json:"prospect_id" example:"P12345"`                                   // Unique prospect ID
	ApplicantName         string         `bson:"applicant_name" json:"applicant_name" example:"John Doe"`                           // Name of the applicant
	MobileNumber          string         `bson:"mobile_number" json:"mobile_number" example:"9876543210"`                           // Mobile number of the applicant
	Gender                string         `bson:"gender" json:"gender" example:"Male"`                                               // Gender of the applicant
	Age                   int            `bson:"age" json:"age" example:"30"`                                                       // Age of the applicant
	ResidentialAddress    string         `bson:"residential_address" json:"residential_address" example:"123 Main Street"`          // Residential address
	YearsOfStay           int            `bson:"years_of_stay" json:"years_of_stay" example:"5"`                                    // Years of stay at the current address
	NumberOfFamilyMembers int            `bson:"number_of_family_members" json:"number_of_family_members" example:"4"`              // Number of family members
	ReferenceName         string         `bson:"reference_name" json:"reference_name" example:"Jane Doe"`                           // Reference name
	ReferenceRelation     string         `bson:"reference_relation" json:"reference_relation" example:"Sister"`                     // Relation with the reference
	ReferenceMobile       string         `bson:"reference_mobile" json:"reference_mobile" example:"9876543211"`                     // Mobile number of the reference
	EmploymentType        string         `bson:"employment_type" json:"employment_type" example:"Employee"`                         // Employment type ("Employee" or "Business")
	OfficeAddress         string         `bson:"office_address" json:"office_address" example:"456 Office Street"`                  // Office address
	YearsInCurrentOffice  int            `bson:"years_in_current_office" json:"years_in_current_office" example:"3"`                // Years in the current office
	Role                  string         `bson:"role" json:"role" example:"Manager"`                                                // Role in the organization
	EmpId                 string         `bson:"emp_id" json:"emp_id" example:"EMP123"`                                             // Employee ID
	Status                ProspectStatus `bson:"status" json:"status" example:"Pending"`                                            // Current status of the prospect
	PreviousExperience    string         `bson:"previous_experience" json:"previous_experience" example:"5 years in sales"`         // Previous experience
	GrossSalary           float64        `bson:"gross_salary" json:"gross_salary" example:"50000.00"`                               // Gross salary
	NetSalary             float64        `bson:"net_salary" json:"net_salary" example:"40000.00"`                                   // Net salary
	ColleagueName         string         `bson:"colleague_name" json:"colleague_name" example:"Mark Smith"`                         // Name of a colleague
	ColleagueDesignation  string         `bson:"colleague_designation" json:"colleague_designation" example:"Team Lead"`            // Designation of the colleague
	ColleagueMobile       string         `bson:"colleague_mobile" json:"colleague_mobile" example:"9876543212"`                     // Mobile number of the colleague
	UploadedImages        []string       `bson:"uploaded_images" json:"uploaded_images" example:"[\"image1.jpg\", \"image2.jpg\"]"` // Uploaded images
	Remarks               string         `bson:"remarks" json:"remarks" example:"Prospect is under review"`                         // Additional remarks
}
