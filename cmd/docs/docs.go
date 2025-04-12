// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/prospects": {
            "post": {
                "description": "Create a new prospect in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Prospects"
                ],
                "summary": "Create a new prospect",
                "parameters": [
                    {
                        "description": "Prospect data",
                        "name": "prospect",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ProspectModel"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controllers.SimpleResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.SimpleResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.SimpleResponse"
                        }
                    }
                }
            }
        },
        "/prospects/{id}": {
            "get": {
                "description": "Retrieve a prospect by their unique ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Prospects"
                ],
                "summary": "Get a prospect by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Prospect ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ProspectModel"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.SimpleResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.SimpleResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Create a new user in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserModel"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controllers.SimpleResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.SimpleResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.SimpleResponse"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Retrieve a user by their unique ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get a user by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserModel"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.SimpleResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.SimpleResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.SimpleResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.ProspectModel": {
            "description": "Prospect model containing all prospect-related information.",
            "type": "object",
            "properties": {
                "age": {
                    "description": "Age of the applicant",
                    "type": "integer",
                    "example": 30
                },
                "applicant_name": {
                    "description": "Name of the applicant",
                    "type": "string",
                    "example": "John Doe"
                },
                "colleague_designation": {
                    "description": "Designation of the colleague",
                    "type": "string",
                    "example": "Team Lead"
                },
                "colleague_mobile": {
                    "description": "Mobile number of the colleague",
                    "type": "string",
                    "example": "9876543212"
                },
                "colleague_name": {
                    "description": "Name of a colleague",
                    "type": "string",
                    "example": "Mark Smith"
                },
                "emp_id": {
                    "description": "Employee ID",
                    "type": "string",
                    "example": "EMP123"
                },
                "employment_type": {
                    "description": "Employment type (\"Employee\" or \"Business\")",
                    "type": "string",
                    "example": "Employee"
                },
                "gender": {
                    "description": "Gender of the applicant",
                    "type": "string",
                    "example": "Male"
                },
                "gross_salary": {
                    "description": "Gross salary",
                    "type": "number",
                    "example": 50000
                },
                "id": {
                    "description": "Incremental ID",
                    "type": "integer",
                    "example": 1
                },
                "mobile_number": {
                    "description": "Mobile number of the applicant",
                    "type": "string",
                    "example": "9876543210"
                },
                "net_salary": {
                    "description": "Net salary",
                    "type": "number",
                    "example": 40000
                },
                "number_of_family_members": {
                    "description": "Number of family members",
                    "type": "integer",
                    "example": 4
                },
                "office_address": {
                    "description": "Office address",
                    "type": "string",
                    "example": "456 Office Street"
                },
                "previous_experience": {
                    "description": "Previous experience",
                    "type": "string",
                    "example": "5 years in sales"
                },
                "prospect_id": {
                    "description": "Unique prospect ID",
                    "type": "string",
                    "example": "P12345"
                },
                "reference_mobile": {
                    "description": "Mobile number of the reference",
                    "type": "string",
                    "example": "9876543211"
                },
                "reference_name": {
                    "description": "Reference name",
                    "type": "string",
                    "example": "Jane Doe"
                },
                "reference_relation": {
                    "description": "Relation with the reference",
                    "type": "string",
                    "example": "Sister"
                },
                "remarks": {
                    "description": "Additional remarks",
                    "type": "string",
                    "example": "Prospect is under review"
                },
                "residential_address": {
                    "description": "Residential address",
                    "type": "string",
                    "example": "123 Main Street"
                },
                "role": {
                    "description": "Role in the organization",
                    "type": "string",
                    "example": "Manager"
                },
                "status": {
                    "description": "Current status of the prospect",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.ProspectStatus"
                        }
                    ],
                    "example": "Pending"
                },
                "uploaded_images": {
                    "description": "Uploaded images",
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "[\"image1.jpg\"",
                        " \"image2.jpg\"]"
                    ]
                },
                "years_in_current_office": {
                    "description": "Years in the current office",
                    "type": "integer",
                    "example": 3
                },
                "years_of_stay": {
                    "description": "Years of stay at the current address",
                    "type": "integer",
                    "example": 5
                }
            }
        },
        "models.ProspectStatus": {
            "type": "string",
            "enum": [
                "Pending",
                "OnVisit",
                "Progress",
                "Approved",
                "Rejected",
                "UnderReview",
                "Completed",
                "Submitted",
                "Cancelled",
                "RePending",
                "Postponed"
            ],
            "x-enum-varnames": [
                "Pending",
                "OnVisit",
                "Progressve",
                "Approved",
                "Rejected",
                "UnderReview",
                "Completed",
                "Submitted",
                "Cancelled",
                "RePending",
                "Postponed"
            ]
        },
        "models.Role": {
            "type": "string",
            "enum": [
                "Admin",
                "Operations Lead",
                "Field Lead",
                "Field Executive",
                "Owner",
                "Operations Executive"
            ],
            "x-enum-varnames": [
                "Admin",
                "OperationsLead",
                "FieldLead",
                "FieldExecutive",
                "Owner",
                "OperationsExecutive"
            ]
        },
        "models.UpdateHistory": {
            "description": "History of updates made to a user.",
            "type": "object",
            "properties": {
                "updated_comments": {
                    "description": "Comments about the update",
                    "type": "string",
                    "example": "Updated user role"
                },
                "updated_time": {
                    "description": "Time of the update",
                    "type": "string",
                    "example": "2023-04-12T15:04:05Z"
                }
            }
        },
        "models.UserModel": {
            "description": "User model containing all user-related information.",
            "type": "object",
            "properties": {
                "created_time": {
                    "description": "Time when the user was created",
                    "type": "string",
                    "example": "2023-04-12T15:04:05Z"
                },
                "password": {
                    "description": "Hashed password",
                    "type": "string",
                    "example": "hashed_password"
                },
                "remarks": {
                    "description": "Additional remarks about the user",
                    "type": "string",
                    "example": "User is active and verified"
                },
                "role": {
                    "description": "Role of the user",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.Role"
                        }
                    ],
                    "example": "Admin"
                },
                "status": {
                    "description": "Status of the user",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.UserStatus"
                        }
                    ],
                    "example": "Active"
                },
                "update_history": {
                    "description": "History of updates",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.UpdateHistory"
                    }
                },
                "updated_time": {
                    "description": "Time when the user was last updated",
                    "type": "string",
                    "example": "2023-04-12T15:04:05Z"
                },
                "user_id": {
                    "description": "Incremental ID",
                    "type": "integer",
                    "example": 1
                },
                "username": {
                    "description": "Username of the user",
                    "type": "string",
                    "example": "john_doe"
                }
            }
        },
        "models.UserStatus": {
            "type": "string",
            "enum": [
                "Created",
                "Confirmed",
                "Verified",
                "Active",
                "Inactive",
                "Disabled",
                "Banned"
            ],
            "x-enum-varnames": [
                "Created",
                "Confirmed",
                "Verified",
                "Active",
                "InActive",
                "Disabled",
                "Banned"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9000",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Kowtha API",
	Description:      "This is the API documentation for the Kowtha backend.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
