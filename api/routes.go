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
		Queries("EndTime", "{EndTime}", "StartTime", "{StartTime}").Methods(http.MethodGet)
	api.HandleFunc("/{account}/spaces/{space}", s.SpaceGetHandler).Methods(http.MethodGet)

	api.HandleFunc("/{account}/instances/{id}/metrics/{metric}.png", s.MetricsGetHandler).Methods(http.MethodGet)
}
