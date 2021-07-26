package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/cost-api/resourcegroupstaggingapi"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// SpaceInventoryGetHandler handles getting the inventory for resources with a spaceid
func (s *server) SpaceInventoryGetHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	spaceID := vars["space"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)
	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
		"arn:aws:iam::aws:policy/AWSResourceGroupsReadOnlyAccess",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	orch := newInventoryOrchestrator(
		resourcegroupstaggingapi.New(resourcegroupstaggingapi.WithSession(session.Session)),
		s.org,
	)

	out, err := orch.GetResourceInventory(r.Context(), account, spaceID)
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(out)
	if err != nil {
		log.Errorf("cannot marshal response (%v) into JSON: %s", out, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Items", strconv.Itoa(len(out)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
