package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog/log"
	"github.com/serdarkalayci/goboiler/webapi/dto"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// APIContext handler for getting and updating Ratings
type APIContext struct {
	v *dto.Validation
}

// DBContext is the struct that has a MongoDB connection together with standard APIContext. It's used for handler functions which will use database
type DBContext struct {
	MongoClient  mongo.Client
	DatabaseName string
	APIContext
}

// NewAPIContext returns a new APIContext handler with the given logger
func NewAPIContext(v *dto.Validation) *APIContext {
	return &APIContext{v}
}

// NewDBContext returns a new DBContext handler with the given logger
func NewDBContext(v *dto.Validation) *DBContext {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// We try to get connectionstring value from the environment variables, if not found it falls back to local database
	connectionString := os.Getenv("ConnectionString")
	if connectionString == "" {
		connectionString = "mongodb://localhost:27017"
		log.Info().Msg("ConnectionString from Env not found, falling back to local DB")
	} else {
		log.Info().Msgf("ConnectionString from Env is used: '%s'", connectionString)
	}
	databaseName := os.Getenv("DatabaseName")
	if databaseName == "" {
		databaseName = "goboiler"
		log.Info().Msg("DatabaseName from Env not found, falling back to default")
	} else {
		log.Info().Msgf("DatabaseName from Env is used: '%s'", databaseName)
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	err = client.Connect(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("An error occured while connecting to tha database")
	} else {

		// Check the connection
		err = client.Ping(context.TODO(), nil)

		if err != nil {
			log.Error().Err(err).Msg("An error occured while connecting to tha database")
		}
		log.Info().Msg("Connected to MongoDB!")
	}
	return &DBContext{*client, databaseName, APIContext{v}}
}

// createSpan creates a new openTracing.Span with the given name and returns it
func createSpan(spanName string, r *http.Request) (span opentracing.Span) {
	tracer := opentracing.GlobalTracer()

	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		// The method is called without a span context in the http header.
		//
		span = tracer.StartSpan(spanName)
	} else {
		// Create the span referring to the RPC client if available.
		// If wireContext == nil, a root span will be created.
		span = opentracing.StartSpan(spanName, ext.RPCServerOption(wireContext))
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, r.URL.RequestURI())
	ext.HTTPMethod.Set(span, r.Method)
	return span
}

// ErrInvalidRatingPath is an error message when the Rating path is not valid
var ErrInvalidRatingPath = fmt.Errorf("Invalid Path, path should be /Details/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}
