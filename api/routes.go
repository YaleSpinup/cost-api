package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *server) routes() {
	api := s.router.PathPrefix("/v1/cost").Subrouter()
	api.HandleFunc("/ping", s.PingHandler).Methods(http.MethodGet)
	api.HandleFunc("/version", s.VersionHandler).Methods(http.MethodGet)
	api.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	api.HandleFunc("/{account}/spaces/{space}", s.SpaceGetHandler).
		Queries("start", "{start}", "end", "{end}").Methods(http.MethodGet)
	api.HandleFunc("/{account}/spaces/{space}", s.SpaceGetHandler).Methods(http.MethodGet)

	api.HandleFunc("/{account}/instances/{id}/metrics/graph.png", s.MetricsGetImageHandler).
		Queries("period", "{period}", "start", "{start}", "end", "{end}").Methods(http.MethodGet)
	api.HandleFunc("/{account}/instances/{id}/metrics/graph.png", s.MetricsGetImageHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/instances/{id}/metrics/graph", s.MetricsGetImageUrlHandler).
		Queries("period", "{period}", "start", "{start}", "end", "{end}").Methods(http.MethodGet)
	api.HandleFunc("/{account}/instances/{id}/metrics/graph", s.MetricsGetImageUrlHandler).Methods(http.MethodGet)

}
