package api

import (
	"fmt"
	"net/http"

	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MetricsGetHandler gets metrics from cloudwatch
func (s *server) MetricsGetHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)

	// get vars from the API route
	account := vars["account"]
	metric := vars["metric"]
	id := vars["id"]

	cwService, ok := s.cloudwatchServices[account]
	if !ok {
		msg := fmt.Sprintf("cloudwatch service not found for account: %s", account)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}
	log.Debugf("found cloudwatch service %+v", cwService)

	out, err := cwService.GetMetricWidget(r.Context(), metric, id)
	if err != nil {
		log.Errorf("failed getting metrics widget image: %s", err)
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
