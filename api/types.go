package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/budgets"
)

// BudgetCreateRequest is the request object to create a Budget
type BudgetCreateRequest struct {
	// Amount in USD for the budget
	Amount string

	// DAILY, MONTHLY, QUARTERLY, or ANNUALLY
	TimeUnit string

	// Alerts is a list of threshold/notification configurations for
	// a budget.  Maximum number is 5.
	Alerts []*BudgetAlert
}

type BudgetAlert struct {
	// The comparison that is used for this notification.
	ComparisonOperator string

	// Whether this notification is in alarm. If a budget notification is in the
	// ALARM state, you have passed the set threshold for the budget.
	NotificationState string

	// Whether the notification is for how much you have spent (ACTUAL) or for how
	// much you're forecasted to spend (FORECASTED).
	NotificationType string

	// The threshold that is associated with a notification. Thresholds are always
	// a percentage, and many customers find value being alerted between 50% - 200%
	// of the budgeted amount. The maximum limit for your threshold is 1,000,000%
	// above the budgeted amount.
	Threshold float64

	// The type of threshold for a notification. For ABSOLUTE_VALUE thresholds,
	// AWS notifies you when you go over or are forecasted to go over your total
	// cost threshold. For PERCENTAGE thresholds, AWS notifies you when you go over
	// or are forecasted to go over a certain percentage of your forecasted spend.
	// For example, if you have a budget for 200 dollars and you have a PERCENTAGE
	// threshold of 80%, AWS notifies you when you go over 160 dollars.
	ThresholdType string

	// Addresses are the email addresses for notifications (up to 10)
	Addresses []string
}

// BudgetResponse is the standard respoonse for a Budget
type BudgetResponse struct {
	Amount   string
	Name     string
	TimeUnit string
	Alerts   []*BudgetAlert
}

func toBudgetAlert(notification *budgets.Notification, subscribers []*budgets.Subscriber) *BudgetAlert {
	addresses := make([]string, len(subscribers))
	for i, s := range subscribers {
		addresses[i] = aws.StringValue(s.Address)
	}

	return &BudgetAlert{
		ComparisonOperator: aws.StringValue(notification.ComparisonOperator),
		NotificationState:  aws.StringValue(notification.NotificationState),
		NotificationType:   aws.StringValue(notification.NotificationType),
		Threshold:          aws.Float64Value(notification.Threshold),
		ThresholdType:      aws.StringValue(notification.ThresholdType),
		Addresses:          addresses,
	}
}

func toBudgetResponse(budget *budgets.Budget, alerts []*BudgetAlert) *BudgetResponse {
	return &BudgetResponse{
		Amount:   aws.StringValue(budget.BudgetLimit.Amount),
		Name:     aws.StringValue(budget.BudgetName),
		TimeUnit: aws.StringValue(budget.TimeUnit),
		Alerts:   alerts,
	}
}
