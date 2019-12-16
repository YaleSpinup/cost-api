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

	out, err := c.GetMetricWidget(context.TODO(), "CPUUtilization", "i-abc12345")
	if err != nil {
		t.Errorf("expected nil error, got: %s", err)
	}

	if !bytes.Equal(out, expected) {
		t.Error("didn't get expected image output from GetMetricWidget")
	}

	// test empty metric
	_, err = c.GetMetricWidget(context.TODO(), "", "i-abc12345")
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrBadRequest {
			t.Errorf("expected error code %s, got: %s", apierror.ErrBadRequest, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}

	// test empty id
	_, err = c.GetMetricWidget(context.TODO(), "CPUUtilization", "")
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrBadRequest {
			t.Errorf("expected error code %s, got: %s", apierror.ErrBadRequest, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}
}
