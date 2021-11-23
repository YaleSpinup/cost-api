package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YaleSpinup/apierror"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

// SpaceGetHandler gets the cost for a space, grouped by the service.  By default,
// it pulls data from the start of the month until now.
func (s *server) SpaceGetHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := s.mapAccountNumber(vars["account"])
	startTime := vars["start"]
	endTime := vars["end"]
	spaceID := vars["space"]
	groupBy := vars["groupby"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)
	policy, err := costExplorerReadPolicy()
	if err != nil {
		handleError(w, apierror.New(apierror.ErrInternalError, "failed to generate policy", err))
		return
	}

	orch, err := s.newCostExplorerOrchestrator(r.Context(), &sessionParams{
		inlinePolicy: policy,
		role:         role,
	})

	out, cached, expire, err := orch.getCostAndUsageForSpace(
		r.Context(),
		&costAndUsageReq{
			account: account,
			spaceID: spaceID,
			start:   startTime,
			end:     endTime,
			groupBy: groupBy,
		},
	)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("X-Cache-Hit", fmt.Sprintf("%t", cached))
	if cached {
		w.Header().Set("X-Cache-Expire", fmt.Sprintf("%0.fs", expire.Seconds()))
	}

	j, err := json.Marshal(out)
	if err != nil {
		log.Errorf("cannot marshal response (%v) into JSON: %s", out, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
