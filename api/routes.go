package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *server) routes() {

	// costs subrouter - /v1/cost
	api := s.router.PathPrefix("/v1/cost").Subrouter()
	api.HandleFunc("/ping", s.PingHandler).Methods(http.MethodGet)
	api.HandleFunc("/version", s.VersionHandler).Methods(http.MethodGet)
	api.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// cost endpoints for a space
	api.HandleFunc("/{account}/spaces/{space}", s.SpaceGetHandler).
		Queries("start", "{start}", "end", "{end}").Methods(http.MethodGet)
	api.HandleFunc("/{account}/spaces/{space}", s.SpaceGetHandler).Methods(http.MethodGet)

	// metrics endpoints for EC2 instances
	// TODO: deprecated but left for backwards compatability, remove me once the UI is updated
	api.HandleFunc("/{account}/instances/{id}/metrics/graph.png", s.GetEC2MetricsImageHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/instances/{id}/metrics/graph", s.GetEC2MetricsURLHandler).Methods(http.MethodGet)

	// metrics subrouter - /v1/metrics
	metricsApi := s.router.PathPrefix("/v1/metrics").Subrouter()
	metricsApi.HandleFunc("/ping", s.PingHandler).Methods(http.MethodGet)
	metricsApi.HandleFunc("/version", s.VersionHandler).Methods(http.MethodGet)
	metricsApi.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// metrics endpoints for EC2 instances
	metricsApi.HandleFunc("/{account}/instances/{id}/graph.png", s.GetEC2MetricsImageHandler).Methods(http.MethodGet)
	metricsApi.HandleFunc("/{account}/instances/{id}/graph", s.GetEC2MetricsURLHandler).Methods(http.MethodGet)

	// metrics endpoints for ECS services
	metricsApi.HandleFunc("/{account}/clusters/{cluster}/services/{service}/graph.png", s.GetECSMetricsImageHandler).Methods(http.MethodGet)
	metricsApi.HandleFunc("/{account}/clusters/{cluster}/services/{service}/graph", s.GetECSMetricsURLHandler).Methods(http.MethodGet)
}
