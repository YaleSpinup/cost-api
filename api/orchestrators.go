package api

import (
	"context"

	"github.com/YaleSpinup/cost-api/budgets"
	"github.com/YaleSpinup/cost-api/computeoptimizer"
	"github.com/YaleSpinup/cost-api/costexplorer"
	"github.com/YaleSpinup/cost-api/resourcegroupstaggingapi"
	"github.com/YaleSpinup/cost-api/sns"
	log "github.com/sirupsen/logrus"
)

type sessionParams struct {
	role         string
	inlinePolicy string
	policyArns   []string
}

type budgetsOrchestrator struct {
	client    *budgets.Budgets
	snsClient *sns.SNS
	org       string
}

type costExplorerOrchestrator struct {
	client *costexplorer.CostExplorer
	server *server
}

type optimizerOrchestrator struct {
	client *computeoptimizer.ComputeOptimizer
	org    string
}

type inventoryOrchestrator struct {
	client *resourcegroupstaggingapi.ResourceGroupsTaggingAPI
	org    string
}

func newBudgetsOrchestrator(budgetsClient *budgets.Budgets, snsClient *sns.SNS, org string) *budgetsOrchestrator {
	return &budgetsOrchestrator{
		client:    budgetsClient,
		snsClient: snsClient,
		org:       org,
	}
}

func newOptimizerOrchestrator(client *computeoptimizer.ComputeOptimizer, org string) *optimizerOrchestrator {
	return &optimizerOrchestrator{
		client: client,
		org:    org,
	}
}

func newInventoryOrchestrator(client *resourcegroupstaggingapi.ResourceGroupsTaggingAPI, org string) *inventoryOrchestrator {
	return &inventoryOrchestrator{
		client: client,
		org:    org,
	}
}

func (s *server) newCostExplorerOrchestrator(ctx context.Context, sp *sessionParams) (*costExplorerOrchestrator, error) {
	log.Debugf("initializing costExplorerOrchestrator")

	session, err := s.assumeRole(
		ctx,
		s.session.ExternalID,
		sp.role,
		sp.inlinePolicy,
		sp.policyArns...,
	)
	if err != nil {
		return nil, err
	}

	return &costExplorerOrchestrator{
		client: costexplorer.New(costexplorer.WithSession(session.Session)),
		server: s,
	}, nil
}
