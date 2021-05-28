package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/budgets"
	log "github.com/sirupsen/logrus"
)

func (o *budgetsOrchestrator) CreateBudget(ctx context.Context, account, spaceID string, req *BudgetCreateRequest) (*BudgetResponse, error) {
	if req.Amount == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "Amount is required", nil)
	}

	if req.TimeUnit == "" {
		req.TimeUnit = "MONTHLY"
	}

	if !validTimeUnit(req.TimeUnit) {
		msg := fmt.Sprintf("invalid time unit %s", req.TimeUnit)
		return nil, apierror.New(apierror.ErrBadRequest, msg, nil)
	}

	budgetName := fmt.Sprintf("%s-%s-01", spaceID, req.TimeUnit)
	spaceTagFilterValue := fmt.Sprintf("user:spinup:spaceid$%s", spaceID)

	// set some reasonable defaults for budgets
	budget := budgets.Budget{
		BudgetName: aws.String(budgetName),
		BudgetLimit: &budgets.Spend{
			Amount: aws.String(req.Amount),
			Unit:   aws.String("USD"),
		},
		BudgetType: aws.String("COST"),
		CostFilters: map[string][]*string{
			"TagKeyValue": {
				aws.String(spaceTagFilterValue),
			},
		},
		CostTypes: &budgets.CostTypes{
			IncludeCredit:            aws.Bool(false),
			IncludeDiscount:          aws.Bool(false),
			IncludeOtherSubscription: aws.Bool(true),
			IncludeRecurring:         aws.Bool(true),
			IncludeRefund:            aws.Bool(false),
			IncludeSubscription:      aws.Bool(true),
			IncludeSupport:           aws.Bool(false),
			IncludeTax:               aws.Bool(false),
			IncludeUpfront:           aws.Bool(false),
			UseAmortized:             aws.Bool(false),
			UseBlended:               aws.Bool(false),
		},
		TimeUnit: aws.String(req.TimeUnit),
	}

	if len(req.Alerts) == 0 {
		return nil, apierror.New(apierror.ErrBadRequest, "at least 1 Alert is required", nil)
	} else if len(req.Alerts) > 5 {
		return nil, apierror.New(apierror.ErrBadRequest, "up to 5 Alerts per budget are supported", nil)
	}

	notifications := []*budgets.NotificationWithSubscribers{}
	for _, a := range req.Alerts {
		log.Debugf("processing alert %+v", a)

		if len(a.Addresses) == 0 {
			return nil, apierror.New(apierror.ErrBadRequest, "at least 1 email address is required per Alert", nil)
		} else if len(a.Addresses) > 10 {
			return nil, apierror.New(apierror.ErrBadRequest, "up to 10 email addresses per Alert are supported", nil)
		}

		subscribers := make([]*budgets.Subscriber, len(a.Addresses))
		for i, s := range a.Addresses {
			subscribers[i] = &budgets.Subscriber{
				Address:          aws.String(s),
				SubscriptionType: aws.String("EMAIL"),
			}
		}

		if !validComparisonOperator(a.ComparisonOperator) {
			msg := fmt.Sprintf("invalid comparison operator '%s', valid values %s", a.NotificationType, strings.Join(budgets.ComparisonOperator_Values(), ", "))
			return nil, apierror.New(apierror.ErrBadRequest, msg, nil)
		}

		if !validNotificationType(a.NotificationType) {
			msg := fmt.Sprintf("invalid notification type '%s', valid values %s", a.NotificationType, strings.Join(budgets.NotificationType_Values(), ", "))
			return nil, apierror.New(apierror.ErrBadRequest, msg, nil)
		}

		if !validThresholdType(a.ThresholdType) {
			msg := fmt.Sprintf("invalid threshold type '%s', valid values %s", a.ThresholdType, strings.Join(budgets.ThresholdType_Values(), ", "))
			return nil, apierror.New(apierror.ErrBadRequest, msg, nil)
		}

		notification := &budgets.Notification{
			ComparisonOperator: aws.String(a.ComparisonOperator),
			NotificationType:   aws.String(a.NotificationType),
			Threshold:          aws.Float64(a.Threshold),
			ThresholdType:      aws.String(a.ThresholdType),
		}

		notifications = append(notifications, &budgets.NotificationWithSubscribers{
			Notification: notification,
			Subscribers:  subscribers,
		})
	}

	if err := o.client.CreateBudget(ctx, &budgets.CreateBudgetInput{
		AccountId:                    aws.String(account),
		Budget:                       &budget,
		NotificationsWithSubscribers: notifications,
	}); err != nil {
		return nil, err
	}

	return toBudgetResponse(&budget, req.Alerts), nil
}

func (o *budgetsOrchestrator) GetBudget(ctx context.Context, account, spaceID, budget string) (*BudgetResponse, error) {
	if !strings.HasPrefix(budget, spaceID) {
		return nil, apierror.New(apierror.ErrBadRequest, "budget doesn't belong to provided space", nil)
	}

	budgetOut, err := o.client.DescribeBudget(ctx, account, budget)
	if err != nil {
		return nil, err
	}

	notifications, err := o.client.DescribeNotifications(ctx, account, budget)
	if err != nil {
		return nil, err
	}

	alerts := []*BudgetAlert{}
	for _, n := range notifications {
		sub, err := o.client.DescribeSubscribers(ctx, account, budget, n)
		if err != nil {
			return nil, err
		}

		alerts = append(alerts, toBudgetAlert(n, sub))
	}

	return toBudgetResponse(budgetOut, alerts), nil
}

func (o *budgetsOrchestrator) ListBudgets(ctx context.Context, account, spaceID string) ([]string, error) {
	out, err := o.client.ListBudgetsWithPrefix(ctx, account, spaceID)
	if err != nil {
		return nil, err
	}

	blist := []string{}
	for _, b := range out {
		blist = append(blist, aws.StringValue(b.BudgetName))
	}

	return blist, nil
}

func (o *budgetsOrchestrator) DeleteBudget(ctx context.Context, account, spaceID, budget string) error {
	if !strings.HasPrefix(budget, spaceID) {
		return apierror.New(apierror.ErrBadRequest, "budget doesn't belong to provided space", nil)
	}

	if err := o.client.DeleteBudget(ctx, account, budget); err != nil {
		return err
	}

	return nil
}

func validComparisonOperator(co string) bool {
	for _, c := range budgets.ComparisonOperator_Values() {
		if c == co {
			return true
		}
	}
	return false
}

func validNotificationType(nt string) bool {
	for _, n := range budgets.NotificationType_Values() {
		if n == nt {
			return true
		}
	}
	return false
}

func validThresholdType(tt string) bool {
	for _, t := range budgets.ThresholdType_Values() {
		if t == tt {
			return true
		}
	}
	return false
}

func validTimeUnit(tu string) bool {
	for _, t := range budgets.TimeUnit_Values() {
		if t == tu {
			return true
		}
	}
	return false
}
