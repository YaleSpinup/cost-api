package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/cost-api/budgets"
	"github.com/YaleSpinup/cost-api/sns"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (s *server) SpaceBudgetsCreatehandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	spaceID := vars["space"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)
	policy, err := budgetReadWritePolicy()
	if err != nil {
		handleError(w, apierror.New(apierror.ErrInternalError, "failed to generate policy", err))
		return
	}

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		policy,
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	req := BudgetCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into create budget input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	orch := newBudgetsOrchestrator(
		budgets.New(budgets.WithSession(session.Session)),
		sns.New(sns.WithSession(session.Session)),
		s.org,
	)

	out, err := orch.CreateBudget(r.Context(), account, spaceID, &req)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *server) SpaceBudgetsListHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	spaceID := vars["space"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)
	policy, err := budgetReadWritePolicy()
	if err != nil {
		handleError(w, apierror.New(apierror.ErrInternalError, "failed to generate policy", err))
		return
	}

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		policy,
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	orch := newBudgetsOrchestrator(
		budgets.New(budgets.WithSession(session.Session)),
		nil,
		s.org,
	)

	out, err := orch.ListBudgets(r.Context(), account, spaceID)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *server) SpaceBudgetsShowHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	spaceID := vars["space"]
	budget := vars["budget"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)
	policy, err := budgetReadWritePolicy()
	if err != nil {
		handleError(w, apierror.New(apierror.ErrInternalError, "failed to generate policy", err))
		return
	}

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		policy,
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	orch := newBudgetsOrchestrator(
		budgets.New(budgets.WithSession(session.Session)),
		nil,
		s.org,
	)

	out, err := orch.GetBudget(r.Context(), account, spaceID, budget)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *server) SpaceBudgetsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	spaceID := vars["space"]
	budget := vars["budget"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)
	policy, err := budgetReadWritePolicy()
	if err != nil {
		handleError(w, apierror.New(apierror.ErrInternalError, "failed to generate policy", err))
		return
	}

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		policy,
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	orch := newBudgetsOrchestrator(
		budgets.New(budgets.WithSession(session.Session)),
		sns.New(sns.WithSession(session.Session)),
		s.org,
	)

	if err := orch.DeleteBudget(r.Context(), account, spaceID, budget); err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
