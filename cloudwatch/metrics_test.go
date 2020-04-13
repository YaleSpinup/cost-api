package cloudwatch

import (
	"bytes"
	"context"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

var testImg = "../img/example_response.png"

func (m *mockCloudwatchClient) GetMetricWidgetImageWithContext(ctx context.Context, input *cloudwatch.GetMetricWidgetImageInput, opts ...request.Option) (*cloudwatch.GetMetricWidgetImageOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	file, err := ioutil.ReadFile(testImg)
	if err != nil {
		return nil, err
	}

	return &cloudwatch.GetMetricWidgetImageOutput{MetricWidgetImage: file}, nil
}

func TestGetMetricWidget(t *testing.T) {
	c := Cloudwatch{
		Service: newmockCloudwatchClient(t, nil),
	}

	expected, err := ioutil.ReadFile(testImg)
	if err != nil {
		t.Errorf("expected nil error reading, got: %s", err)
	}

	req := MetricsRequest{
		"metrics": []Metric{
			{"AWS/EC2", "CPUUtilization", "InstanceId", "i-abc12345"},
		},
		"period": int64(300),
		"start":  "-P1D",
		"end":    "PT0H",
	}

	out, err := c.GetMetricWidget(context.TODO(), req)
	if err != nil {
		t.Errorf("expected nil error, got: %s", err)
	}

	if !bytes.Equal(out, expected) {
		t.Error("didn't get expected image output from GetMetricWidget")
	}

	// test nil metric request
	_, err = c.GetMetricWidget(context.TODO(), nil)
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrBadRequest {
			t.Errorf("expected error code %s, got: %s", apierror.ErrBadRequest, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}
}
