package api

import (
	"github.com/YaleSpinup/cost-api/budgets"
	"github.com/YaleSpinup/cost-api/computeoptimizer"
	"github.com/YaleSpinup/cost-api/sns"
)

type budgetsOrchestrator struct {
	client    *budgets.Budgets
	snsClient *sns.SNS
	org       string
}

func newBudgetsOrchestrator(budgetsClient *budgets.Budgets, snsClient *sns.SNS, org string) *budgetsOrchestrator {
	return &budgetsOrchestrator{
		client:    budgetsClient,
		snsClient: snsClient,
		org:       org,
	}
}

type optimizerOrchestrator struct {
	client *computeoptimizer.ComputeOptimizer
	org    string
}

func newOptimizerOrchestrator(client *computeoptimizer.ComputeOptimizer, org string) *optimizerOrchestrator {
	return &optimizerOrchestrator{
		client: client,
		org:    org,
	}
}
