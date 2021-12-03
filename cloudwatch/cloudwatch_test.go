package cloudwatch

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

// mockCloudwatchClient is a fake cloudwatch client
type mockCloudwatchClient struct {
	cloudwatchiface.CloudWatchAPI
	t   *testing.T
	err error
}

func newmockCloudwatchClient(t *testing.T, err error) cloudwatchiface.CloudWatchAPI {
	return &mockCloudwatchClient{
		t:   t,
		err: err,
	}
}

func TestNewSession(t *testing.T) {
	e := New()
	if to := reflect.TypeOf(e).String(); to != "*cloudwatch.Cloudwatch" {
		t.Errorf("expected type to be 'cloudwatch.Cloudwatch', got %s", to)
	}
}
