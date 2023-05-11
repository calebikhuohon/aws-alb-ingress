package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	router    http.Handler
	awsRegion string
	s3Bucket  string
	s3Item    string
}

type Config struct {
	Timeout time.Duration
}

func New(
	config Config,
	awsRegion string,
	s3Bucket string,
	s3Item string,
) http.Handler {
	r := chi.NewRouter()

	h := &Handler{
		router:    r,
		awsRegion: awsRegion,
		s3Bucket:  s3Bucket,
		s3Item:    s3Item,
	}

	timeout := 10 * time.Second
	if config.Timeout > 0 {
		timeout = config.Timeout
	}

	r.Use(
		chimiddleware.Timeout(timeout),
		chimiddleware.SetHeader("Content-Type", "application/json"),
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedHeaders:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "HEAD", "OPTIONS", "PUT"},
			Debug:            false,
		}).Handler,
	)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	r.Get("/", h.GetDefault)

	r.Get("/get-res", h.GetResult)

	return h
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h Handler) GetDefault(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{"msg": "ok"})
}

func (h Handler) GetResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := GetResult(ctx, h.awsRegion, h.s3Bucket, h.s3Item)
	if err != nil {
		log.WithError(err).Errorf("failed to get results from s3")
	}
	if err != nil {
		JSONError(w, Error{
			Status:  false,
			Message: "failed to get results",
			Data:    nil,
		}, http.StatusBadRequest)
	}

	_ = json.NewEncoder(w).Encode(result)
}
