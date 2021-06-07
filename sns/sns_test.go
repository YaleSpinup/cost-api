package sns

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// mockSNSClient is a fake ecs client
type mockSNSClient struct {
	snsiface.SNSAPI
	t   *testing.T
	err error
}

func newMockBudgetsClient(t *testing.T, err error) snsiface.SNSAPI {
	return &mockSNSClient{
		t:   t,
		err: err,
	}
}

func TestNewSession(t *testing.T) {
	client := New()
	to := reflect.TypeOf(client).String()
	if to != "*sns.SNS" {
		t.Errorf("expected type to be '*sns.SNS', got %s", to)
	}
}
