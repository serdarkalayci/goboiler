package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	openapimw "github.com/go-openapi/runtime/middleware"

	"github.com/serdarkalayci/goboiler/webapi/dto"
	"github.com/serdarkalayci/goboiler/webapi/handlers"
	"github.com/serdarkalayci/goboiler/webapi/middleware"

	"github.com/rs/zerolog"

	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"github.com/rs/zerolog/log"
	config "github.com/serdarkalayci/goboiler/webapi/configuration"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	jprom "github.com/uber/jaeger-lib/metrics/prometheus"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var bindAddress = env.String("BASE_URL", false, ":5500", "Bind address for the server")

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	config.SetConfigValues()

	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg, err := jaegercfg.FromEnv()
	if err != nil || cfg.ServiceName == "" {
		cfg = &jaegercfg.Configuration{
			ServiceName: "BookInfo.Details",
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans: true,
			},
		}
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := jprom.New()

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, _ := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	env.Parse()

	v := dto.NewValidation()

	// create the handlers
	apiContext := handlers.NewAPIContext(v)
	dbContext := handlers.NewDBContext(v)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()
	sm.Use(middleware.MetricsMiddleware)

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/", apiContext.Index)
	getR.HandleFunc("/health/live", apiContext.Live)
	getR.HandleFunc("/health/ready", dbContext.Ready)
	getR.HandleFunc("/products/{id}", dbContext.GetSingleProduct)
	getR.HandleFunc("/products", dbContext.GetAllProducts)

	// handler for documentation
	opts := openapimw.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := openapimw.Redoc(opts, nil)

	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	sm.PathPrefix("/metrics").Handler(promhttp.Handler())
	prometheus.MustRegister(middleware.RequestCounterVec)
	prometheus.MustRegister(middleware.RequestDurationGauge)

	// start the server
	go func() {
		log.Debug().Msgf("Starting server on %s", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			log.Error().Err(err).Msg("Error starting server")
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Info().Msgf("Got signal: %s", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
