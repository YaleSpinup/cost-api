package api

import "github.com/YaleSpinup/cost-api/budgets"

type budgetsOrchestrator struct {
	client budgets.Budgets
	org    string
}

func newBudgetsOrchestrator(client budgets.Budgets, org string) *budgetsOrchestrator {
	return &budgetsOrchestrator{
		client: client,
		org:    org,
	}
}
