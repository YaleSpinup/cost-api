package budgets

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/budgets"
	"github.com/aws/aws-sdk-go/service/budgets/budgetsiface"
)

// mockBudgetsClient is a fake ecs client
type mockBudgetsClient struct {
	budgetsiface.BudgetsAPI
	t   *testing.T
	err error
}

func newMockBudgetsClient(t *testing.T, err error) budgetsiface.BudgetsAPI {
	return &mockBudgetsClient{
		t:   t,
		err: err,
	}
}

func TestNewSession(t *testing.T) {
	client := New()
	to := reflect.TypeOf(client).String()
	if to != "*budgets.Budgets" {
		t.Errorf("expected type to be '*budgets.Budgets', got %s", to)
	}
}

func (m *mockBudgetsClient) CreateBudgetWithContext(ctx context.Context, input *budgets.CreateBudgetInput, opts ...request.Option) (*budgets.CreateBudgetOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &budgets.CreateBudgetOutput{}, nil
}

func TestBudgets_CreateBudget(t *testing.T) {
	type fields struct {
		session *session.Session
		Service budgetsiface.BudgetsAPI
	}
	type args struct {
		ctx   context.Context
		input *budgets.CreateBudgetInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "nil input",
			fields: fields{Service: newMockBudgetsClient(t, nil)},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name:   "nil budget",
			fields: fields{Service: newMockBudgetsClient(t, nil)},
			args: args{
				ctx:   context.TODO(),
				input: &budgets.CreateBudgetInput{},
			},
			wantErr: true,
		},
		{
			name:   "aws err",
			fields: fields{Service: newMockBudgetsClient(t, awserr.New(budgets.ErrCodeCreationLimitExceededException, "boom", nil))},
			args: args{
				ctx: context.TODO(),
				input: &budgets.CreateBudgetInput{
					AccountId: aws.String("0123456789"),
					Budget:    &budgets.Budget{},
				},
			},
			wantErr: true,
		},
		{
			name:   "valid input",
			fields: fields{Service: newMockBudgetsClient(t, nil)},
			args: args{
				ctx: context.TODO(),
				input: &budgets.CreateBudgetInput{
					AccountId: aws.String("0123456789"),
					Budget:    &budgets.Budget{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Budgets{
				session: tt.fields.session,
				Service: tt.fields.Service,
			}
			if err := b.CreateBudget(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Budgets.CreateBudget() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBudgets_ListBudgetsWithPrefix(t *testing.T) {
	type fields struct {
		session *session.Session
		Service budgetsiface.BudgetsAPI
	}
	type args struct {
		ctx     context.Context
		account string
		prefix  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*budgets.Budget
		wantErr bool
	}{
		{
			name:   "empty account",
			fields: fields{Service: newMockBudgetsClient(t, nil)},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Budgets{
				session: tt.fields.session,
				Service: tt.fields.Service,
			}
			got, err := b.ListBudgetsWithPrefix(tt.args.ctx, tt.args.account, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Budgets.ListBudgetsWithPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Budgets.ListBudgetsWithPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBudgets_DeleteBudget(t *testing.T) {
	type fields struct {
		session *session.Session
		Service budgetsiface.BudgetsAPI
	}
	type args struct {
		ctx     context.Context
		account string
		budget  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "empty input",
			fields: fields{Service: newMockBudgetsClient(t, nil)},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Budgets{
				session: tt.fields.session,
				Service: tt.fields.Service,
			}
			if err := b.DeleteBudget(tt.args.ctx, tt.args.account, tt.args.budget); (err != nil) != tt.wantErr {
				t.Errorf("Budgets.DeleteBudget() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBudgets_DescribeBudget(t *testing.T) {
	type fields struct {
		session *session.Session
		Service budgetsiface.BudgetsAPI
	}
	type args struct {
		ctx     context.Context
		account string
		budget  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *budgets.Budget
		wantErr bool
	}{
		{
			name:   "empty input",
			fields: fields{Service: newMockBudgetsClient(t, nil)},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Budgets{
				session: tt.fields.session,
				Service: tt.fields.Service,
			}
			got, err := b.DescribeBudget(tt.args.ctx, tt.args.account, tt.args.budget)
			if (err != nil) != tt.wantErr {
				t.Errorf("Budgets.DescribeBudget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Budgets.DescribeBudget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBudgets_DescribeNotifications(t *testing.T) {
	type fields struct {
		session *session.Session
		Service budgetsiface.BudgetsAPI
	}
	type args struct {
		ctx     context.Context
		account string
		budget  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*budgets.Notification
		wantErr bool
	}{
		{
			name:   "empty input",
			fields: fields{Service: newMockBudgetsClient(t, nil)},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Budgets{
				session: tt.fields.session,
				Service: tt.fields.Service,
			}
			got, err := b.DescribeNotifications(tt.args.ctx, tt.args.account, tt.args.budget)
			if (err != nil) != tt.wantErr {
				t.Errorf("Budgets.DescribeNotifications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Budgets.DescribeNotifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBudgets_DescribeSubscribers(t *testing.T) {
	type fields struct {
		session *session.Session
		Service budgetsiface.BudgetsAPI
	}
	type args struct {
		ctx          context.Context
		account      string
		budget       string
		notification *budgets.Notification
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*budgets.Subscriber
		wantErr bool
	}{
		{
			name:   "empty input",
			fields: fields{Service: newMockBudgetsClient(t, nil)},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Budgets{
				session: tt.fields.session,
				Service: tt.fields.Service,
			}
			got, err := b.DescribeSubscribers(tt.args.ctx, tt.args.account, tt.args.budget, tt.args.notification)
			if (err != nil) != tt.wantErr {
				t.Errorf("Budgets.DescribeSubscribers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Budgets.DescribeSubscribers() = %v, want %v", got, tt.want)
			}
		})
	}
}
