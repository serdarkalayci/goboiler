package handlers

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/serdarkalayci/goboiler/webapi/data"
	"github.com/serdarkalayci/goboiler/webapi/dto"
)

// KeyProduct is a key used carrying the Product object within the context, just to avoid deserializing it multiple times
type KeyProduct struct{}

// MiddlewareValidateNewProduct Product new book product in the request and calls next if ok
func (apiContext *APIContext) MiddlewareValidateNewProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := &dto.Product{}

		err := data.FromJSON(product, r.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error deserializing book product")

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		// validate the product
		errs := apiContext.v.Validate(product)
		if len(errs) != 0 {
			log.Error().Err(errs[0]).Msg("Error validating book product")

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add the rating to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

// MiddlewareValidateProductPrice validates new book product in the request and calls next if ok
// func (apiContext *APIContext) MiddlewareValidateProductPrice(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		detprice := &dto.ProductPrice{}

// 		err := data.FromJSON(detprice, r.Body)
// 		if err != nil {
// 			log.Error().Err(err).Msg("Error deserializing price data")

// 			rw.WriteHeader(http.StatusBadRequest)
// 			data.ToJSON(&GenericError{Message: err.Error()}, rw)
// 			return
// 		}

// 		// validate the product
// 		errs := apiContext.v.Validate(detprice)
// 		if len(errs) != 0 {
// 			log.Error().Err(errs[0]).Msg("Error validating book product")

// 			// return the validation messages as an array
// 			rw.WriteHeader(http.StatusUnprocessableEntity)
// 			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
// 			return
// 		}

// 		// add the rating to the context
// 		ctx := context.WithValue(r.Context(), KeyProduct{}, detprice)
// 		r = r.WithContext(ctx)

// 		// Call the next handler, which can be another middleware in the chain, or the final handler.
// 		next.ServeHTTP(rw, r)
// 	})
// }
