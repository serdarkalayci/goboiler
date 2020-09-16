package data

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

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

	// the Feature list of the product
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

// GetProductByID returns a single Product which matches the id from the
// database.
// If a Product is not found this function returns a ProductNotFound error
func GetProductByID(id primitive.ObjectID, dbClient mongo.Client, dbName string) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := dbClient.Database(dbName).Collection("products")
	var product Product
	za := primitive.ObjectID.String(id)
	log.Debug().Msgf("Getting the project from database with id: %s", za)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProducts returns all Products from the database.
func GetProducts(dbClient mongo.Client, dbName string) (*[]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := dbClient.Database(dbName).Collection("products")
	var products []Product
	log.Debug().Msg("Getting all the projects from database")
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var product Product
		err := cur.Decode(&product)
		if err != nil {
			log.Error().Err(err).Msg("Result cannot be decoded into Product")
			return nil, err
		}
		products = append(products, product)
	}
	return &products, nil
}
