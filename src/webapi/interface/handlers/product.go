package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/serdarkalayci/goboiler/webapi/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetSingleProduct gets a single product from database
// swagger:route GET /products/{id} Products getSingleProduct
// Return a list of Product from the database
// responses:
//	200: ProductResponse
//	404: errorResponse
// ListSingle handles GET requests
func (ctx *DBContext) GetSingleProduct(rw http.ResponseWriter, r *http.Request) {
	span := createSpan("api.Product.GetSingleProduct", r)
	defer span.Finish()

	id := getProductID(r)

	log.Debug().Msgf("get record id %d", id)

	product, err := data.GetProductByID(id, ctx.MongoClient, ctx.DatabaseName)
	if err != nil {
		log.Error().Err(err).Msg("Error getting Detail")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	err = data.ToJSON(product, rw)
	if err != nil {
		// we should never be here but log the error just incase
		log.Error().Err(err).Msg("Error serializing product")
	}
}

// GetAllProducts gets all products from database
// swagger:route GET /products Products getAllProducts
// Return a list of Product from the database
// responses:
//	200: ProductsResponse
//	404: errorResponse
// ListSingle handles GET requests
func (ctx *DBContext) GetAllProducts(rw http.ResponseWriter, r *http.Request) {
	span := createSpan("api.Product.GetSingleProduct", r)
	defer span.Finish()

	log.Debug().Msgf("get all products initiated")

	products, err := data.GetProducts(ctx.MongoClient, ctx.DatabaseName)
	if err != nil {
		log.Error().Err(err).Msg("Error getting Product")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	err = data.ToJSON(products, rw)
	if err != nil {
		// we should never be here but log the error just incase
		log.Error().Err(err).Msg("Error serializing products")
	}
}

// getProductID returns the ProductID from the URL
// Panics if cannot convert the id into an bson.primitive.ObjectID
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) primitive.ObjectID {
	// parse the Rating id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
