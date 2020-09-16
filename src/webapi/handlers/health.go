package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/serdarkalayci/goboiler/webapi/data"
)

// swagger:route GET /health/live Health Live
// Return 200 if the api is up and running
// responses:
//	200: OK
//	404: errorResponse

// Live handles GET requests
func (ctx *APIContext) Live(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

// swagger:route GET /health/ready Health Ready
// Return 200 if the api is up and running and connected to the database
// responses:
//	200: OK
//	404: errorResponse

// Ready handles GET requests
func (ctx *DBContext) Ready(rw http.ResponseWriter, r *http.Request) {
	err := data.GetHealth(ctx.MongoClient, ctx.DatabaseName)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to database")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
