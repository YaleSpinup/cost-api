package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *server) routes() {

	// costs subrouter - /v1/cost
	api := s.router.PathPrefix("/v1/cost").Subrouter()
	api.HandleFunc("/ping", s.PingHandler).Methods(http.MethodGet)
	api.HandleFunc("/version", s.VersionHandler).Methods(http.MethodGet)
	api.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// cost endpoints for a space
	api.HandleFunc("/{account}/spaces/{space}", s.SpaceGetHandler).Methods(http.MethodGet).MatcherFunc(matchSpaceQueries)

	api.HandleFunc("/{account}/spaces/{space}/budgets", s.SpaceBudgetsCreatehandler).Methods(http.MethodPost)
	api.HandleFunc("/{account}/spaces/{space}/budgets", s.SpaceBudgetsListHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/spaces/{space}/budgets/{budget}", s.SpaceBudgetsShowHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/spaces/{space}/budgets/{budget}", s.SpaceBudgetsDeleteHandler).Methods(http.MethodDelete)

	api.HandleFunc("/{account}/spaces/{space}/instances/{id}/optimizer", s.SpaceInstanceOptimizer).Methods(http.MethodGet)

	// metrics subrouter - /v1/metrics
	metricsApi := s.router.PathPrefix("/v1/metrics").Subrouter()
	metricsApi.HandleFunc("/ping", s.PingHandler).Methods(http.MethodGet)
	metricsApi.HandleFunc("/version", s.VersionHandler).Methods(http.MethodGet)
	metricsApi.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// metrics endpoints for EC2 instances
	metricsApi.HandleFunc("/{account}/instances/{id}/graph", s.GetEC2MetricsURLHandler).Methods(http.MethodGet)
	// metrics endpoints for ECS services
	metricsApi.HandleFunc("/{account}/clusters/{cluster}/services/{service}/graph", s.GetECSMetricsURLHandler).Methods(http.MethodGet)
	// metrics endpoints for S3 buckets
	metricsApi.HandleFunc("/{account}/buckets/{bucket}/graph", s.GetS3MetricsURLHandler).Queries("metric", "{metric:(?:BucketSizeBytes|NumberOfObjects)}").Methods(http.MethodGet)
	// metrics endpoints for RDS services
	metricsApi.HandleFunc("/{account}/rds/{type}/{id}/graph", s.GetRDSMetricsURLHandler).Methods(http.MethodGet)

	inventoryApi := s.router.PathPrefix("/v1/inventory").Subrouter()
	inventoryApi.HandleFunc("/ping", s.PingHandler).Methods(http.MethodGet)
	inventoryApi.HandleFunc("/version", s.VersionHandler).Methods(http.MethodGet)
	inventoryApi.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// inventory endpoints for a space
	inventoryApi.HandleFunc("/{account}/spaces/{space}", s.SpaceInventoryGetHandler).Methods(http.MethodGet)
}

// custom matcher for space queries
func matchSpaceQueries(req *http.Request, r *mux.RouteMatch) bool {
	queries := req.URL.Query()
	if len(queries) == 0 {
		return true
	}

	if r.Vars == nil {
		r.Vars = make(map[string]string)
	}

	s, sok := queries["start"]
	if sok {
		r.Vars["start"] = s[0]
	}

	e, eok := queries["end"]
	if eok {
		r.Vars["end"] = e[0]
	}

	// start and end must be in the same state
	if sok != eok {
		return false
	}

	if g, ok := queries["groupby"]; ok {
		r.Vars["groupby"] = g[0]
	}

	return true
}
