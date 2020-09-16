package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product defines the structure for a product
// swagger:model
type Product struct {
	// the id of the product
	//
	// required: false
	ID primitive.ObjectID `json:"id" bson:"_id"` // Unique identifier for the book

	// the name of the product
	//
	// required: true
	Name string `json:"name" bson:"name" validate:"required"`

	// the ISBNFeature list of the product
	//
	// required: false
	Features []Feature `json:"features" bson:"features"`
}

// Feature defines the structure for each feature of a product
// swagger:model
type Feature struct {
	// the id of the feature
	//
	// required: false
	// min: 1
	ID primitive.ObjectID `json:"id" bson:"id"` // Unique identifier for the feature

	// the user friendly name of the feature
	//
	// required: true
	Name string `json:"name" bson:"name" validate:"required"`

	// the code friendly name of the feature
	//
	// required: true
	Code string `json:"code" bson:"code" validate:"required"`

	// the code friendly name of the feature
	//
	// required: true
	Type FlagType `json:"type" bson:"type" validate:"required"`
}

// FlagType is the enum that enumerates the type of feature flag
type FlagType int

const (
	// Bool indicates the flag is either true or false
	Bool FlagType = iota
	// Int indicates the flag is an integer
	Int
	// Date indicates the flag is date bound flag
	Date
)
