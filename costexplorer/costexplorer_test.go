package costexplorer

import (
	"reflect"
	"testing"

	"github.com/YaleSpinup/cost-api/common"
	"github.com/aws/aws-sdk-go/service/costexplorer/costexploreriface"
)

// mockCostExplorerClient is a fake costexplorer client
type mockCostExplorerClient struct {
	costexploreriface.CostExplorerAPI
	t   *testing.T
	err error
}

func newmockCostExplorerClient(t *testing.T, err error) costexploreriface.CostExplorerAPI {
	return &mockCostExplorerClient{
		t:   t,
		err: err,
	}
}

func TestNewSession(t *testing.T) {
	e := NewSession(common.Account{})
	if to := reflect.TypeOf(e).String(); to != "costexplorer.CostExplorer" {
		t.Errorf("expected type to be 'costexplorer.CostExplorer', got %s", to)
	}
}
