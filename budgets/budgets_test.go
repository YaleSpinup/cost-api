package budgets

import (
	"reflect"
	"testing"

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
	if to != "budgets.Budgets" {
		t.Errorf("expected type to be 'budgets.Budgets', got %s", to)
	}
}
