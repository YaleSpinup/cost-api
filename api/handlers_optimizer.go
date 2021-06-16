package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/cost-api/computeoptimizer"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (s *server) SpaceInstanceOptimizer(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	instanceID := vars["id"]

	var j []byte
	var err error
	if c, expire, ok := s.optimizerCache.GetWithExpiration(instanceID); ok && c != nil {
		log.Debugf("found optimizer result for %s in the cache, returning", instanceID)

		w.Header().Set("X-Cache-Hit", "true")
		w.Header().Set("X-Cache-Expire", fmt.Sprintf("%0.fs", time.Until(expire).Seconds()))

		if j, err = json.Marshal(c); err != nil {
			log.Errorf("cannot marshal response (%v) into JSON: %s", c, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)
		session, err := s.assumeRole(
			r.Context(),
			s.session.ExternalID,
			role,
			"",
			"arn:aws:iam::aws:policy/ComputeOptimizerReadOnlyAccess",
		)
		if err != nil {
			msg := fmt.Sprintf("failed to assume role in account: %s", account)
			handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
			return
		}

		orch := newOptimizerOrchestrator(
			computeoptimizer.New(computeoptimizer.WithSession(session.Session)),
			s.org,
		)

		out, err := orch.GetInstanceRecommendations(r.Context(), account, instanceID)
		if err != nil {
			handleError(w, err)
			return
		}

		if j, err = json.Marshal(out); err != nil {
			log.Errorf("cannot marshal response (%v) into JSON: %s", out, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// cache results
		s.optimizerCache.SetDefault(instanceID, out)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
