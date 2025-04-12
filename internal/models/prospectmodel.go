package models

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

type ProspectModel struct {
	ID                    int            `bson:"id" json:"ID"` // Incremental ID
	ProspectId            string         `bson:"prospect_id" json:"prospect_id"`
	ApplicantName         string         `bson:"applicant_name" json:"applicant_name"`
	MobileNumber          string         `bson:"mobile_number" json:"mobile_number"`
	Gender                string         `bson:"gender" json:"gender"`
	Age                   int            `bson:"age" json:"age"`
	ResidentialAddress    string         `bson:"residential_address" json:"residential_address"`
	YearsOfStay           int            `bson:"years_of_stay" json:"years_of_stay"`
	NumberOfFamilyMembers int            `bson:"number_of_family_members" json:"number_of_family_members"`
	ReferenceName         string         `bson:"reference_name" json:"reference_name"`
	ReferenceRelation     string         `bson:"reference_relation" json:"reference_relation"`
	ReferenceMobile       string         `bson:"reference_mobile" json:"reference_mobile"`
	EmploymentType        string         `bson:"employment_type" json:"employment_type"` // "Employee" or "Business"
	OfficeAddress         string         `bson:"office_address" json:"office_address"`
	YearsInCurrentOffice  int            `bson:"years_in_current_office" json:"years_in_current_office"`
	Role                  string         `bson:"role" json:"role"`
	EmpId                 string         `bson:"emp_id" json:"emp_id"`
	Status                ProspectStatus `bson:"status" json:"status"`
	PreviousExperience    string         `bson:"previous_experience" json:"previous_experience"`
	GrossSalary           float64        `bson:"gross_salary" json:"gross_salary"`
	NetSalary             float64        `bson:"net_salary" json:"net_salary"`
	ColleagueName         string         `bson:"colleague_name" json:"colleague_name"`
	ColleagueDesignation  string         `bson:"colleague_designation" json:"colleague_designation"`
	ColleagueMobile       string         `bson:"colleague_mobile" json:"colleague_mobile"`
	UploadedImages        []string       `bson:"uploaded_images" json:"uploaded_images"`
	Remarks               string         `bson:"remarks" json:"remarks"`
}
