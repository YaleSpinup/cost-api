package budgets

import (
	"context"
	"strings"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/budgets"
	"github.com/aws/aws-sdk-go/service/budgets/budgetsiface"
	log "github.com/sirupsen/logrus"
)

type Budgets struct {
	session *session.Session
	Service budgetsiface.BudgetsAPI
}

type BudgetsOption func(*Budgets)

func New(opts ...BudgetsOption) *Budgets {
	client := Budgets{}

	for _, opt := range opts {
		opt(&client)
	}

	if client.session != nil {
		client.Service = budgets.New(client.session)
	}

	return &client
}

func WithSession(sess *session.Session) BudgetsOption {
	return func(client *Budgets) {
		log.Debug("using aws session")
		client.session = sess
	}
}

func WithCredentials(key, secret, token, region string) BudgetsOption {
	return func(client *Budgets) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		client.session = sess
	}
}

func (b *Budgets) CreateBudget(ctx context.Context, input *budgets.CreateBudgetInput) error {
	if input == nil || input.Budget == nil {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("creating budget %s", aws.StringValue(input.Budget.BudgetName))

	if _, err := b.Service.CreateBudgetWithContext(ctx, input); err != nil {
		return ErrCode("failed to create Budget", err)
	}

	return nil
}

func (b *Budgets) ListBudgetsWithPrefix(ctx context.Context, account, prefix string) ([]*budgets.Budget, error) {
	if account == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("listing budgets in account %s with prefix %s", account, prefix)

	list := []*budgets.Budget{}
	input := budgets.DescribeBudgetsInput{
		AccountId: aws.String(account),
	}

	for {
		out, err := b.Service.DescribeBudgetsWithContext(ctx, &input)
		if err != nil {
			return nil, ErrCode("failed to describe budgets", err)
		}

		for _, budget := range out.Budgets {
			name := aws.StringValue(budget.BudgetName)
			if strings.HasPrefix(name, prefix) {
				list = append(list, budget)
			}
		}

		if out.NextToken != nil {
			input.NextToken = out.NextToken
			continue
		}

		log.Debugf("returning list of budgets: %+v", list)

		return list, nil
	}
}

func (b *Budgets) DeleteBudget(ctx context.Context, account, budget string) error {
	if account == "" || budget == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("deleting budget %s in account %s", budget, account)

	if _, err := b.Service.DeleteBudgetWithContext(ctx, &budgets.DeleteBudgetInput{
		AccountId:  aws.String(account),
		BudgetName: aws.String(budget),
	}); err != nil {
		return ErrCode("failed to delete budget", err)
	}

	return nil
}

func (b *Budgets) DescribeBudget(ctx context.Context, account, budget string) (*budgets.Budget, error) {
	if account == "" || budget == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("describing budget %s in account %s", budget, account)

	out, err := b.Service.DescribeBudgetWithContext(ctx, &budgets.DescribeBudgetInput{
		AccountId:  aws.String(account),
		BudgetName: aws.String(budget),
	})
	if err != nil {
		return nil, ErrCode("failed to describe budget", err)
	}

	log.Debugf("output describing budget: %+v", out)

	return out.Budget, nil
}

func (b *Budgets) DescribeNotifications(ctx context.Context, account, budget string) ([]*budgets.Notification, error) {
	if account == "" || budget == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("describing notifications for budget %s in account %s", budget, account)

	out, err := b.Service.DescribeNotificationsForBudgetWithContext(ctx, &budgets.DescribeNotificationsForBudgetInput{
		AccountId:  aws.String(account),
		BudgetName: aws.String(budget),
		MaxResults: aws.Int64(100),
	})
	if err != nil {
		return nil, ErrCode("failed to describe budget", err)
	}

	log.Debugf("output describing budget notification: %+v", out)

	return out.Notifications, nil
}

func (b *Budgets) DescribeSubscribers(ctx context.Context, account, budget string, notification *budgets.Notification) ([]*budgets.Subscriber, error) {
	if account == "" || budget == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("describing subscribers for notifications for budget %s in account %s", budget, account)

	out, err := b.Service.DescribeSubscribersForNotificationWithContext(ctx, &budgets.DescribeSubscribersForNotificationInput{
		AccountId:    aws.String(account),
		BudgetName:   aws.String(budget),
		MaxResults:   aws.Int64(100),
		Notification: notification,
	})
	if err != nil {
		return nil, ErrCode("failed to describe budget", err)
	}

	log.Debugf("output describing budget notification subscribers: %+v", out)

	return out.Subscribers, nil
}
