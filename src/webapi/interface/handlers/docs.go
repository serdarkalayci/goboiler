// Package classification of Details API
//
// Documentation for Details API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"github.com/serdarkalayci/goboiler/webapi/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//
// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handers

// Generic error message returned as a string
// swagger:response OK
type okResponseWrapper struct {
}

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Validation errors defined as an array of strings
// swagger:response errorValidation
type errorValidationWrapper struct {
	// Collection of the errors
	// in: body
	Body ValidationError
}

// A list of products
// swagger:response ProductsResponse
type productsResponseWrapper struct {
	// All current products
	// in: body
	Body []data.Product
}

// Data structure representing a single product
// swagger:response ProductResponse
type productResponseWrapper struct {
	// Newly created product
	// in: body
	Body data.Product
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:parameters updateDetail
type productIDParamsWrapper struct {
	// The id of the product for which the operation relates
	// in: path
	// required: true
	ID primitive.ObjectID `json:"id"`
}
