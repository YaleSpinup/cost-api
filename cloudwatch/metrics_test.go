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

func TestMetricRequestString(t *testing.T) {
	type metricRequestString struct {
		input  *MetricsRequest
		output string
	}

	tests := []metricRequestString{
		{
			output: "",
			input:  nil,
		},
		{
			output: "/end:PT0H/height:400/period:300/start:-P1D/stat:Average/width:600",
			input: &MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Average",
				"height": int64(400),
				"width":  int64(600),
			},
		},
		{
			output: "/end:PT0H/height:400/period:10/start:-PT1H/stat:Maximum/width:600",
			input: &MetricsRequest{
				"start":  "-PT1H",
				"end":    "PT0H",
				"period": int64(10),
				"stat":   "Maximum",
				"height": int64(400),
				"width":  int64(600),
			},
		},
		{
			output: "/end:PT0H/height:400/period:300/start:-PT1H/stat:Average/width:600",
			input: &MetricsRequest{
				"start":  "-PT1H",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Average",
				"height": int64(400),
				"width":  int64(600),
			},
		},
		{
			output: "/end:-PT1H/height:400/period:300/start:-P1D/stat:Average/width:600",
			input: &MetricsRequest{
				"start":  "-P1D",
				"end":    "-PT1H",
				"period": int64(300),
				"stat":   "Average",
				"height": int64(400),
				"width":  int64(600),
			},
		},
		{
			output: "/end:PT0H/height:400/period:10/start:-P1D/stat:Average/width:600",
			input: &MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(10),
				"stat":   "Average",
				"height": int64(400),
				"width":  int64(600),
			},
		},
		{
			output: "/end:PT0H/height:400/period:300/start:-P1D/stat:Minimum/width:600",
			input: &MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Minimum",
				"height": int64(400),
				"width":  int64(600),
			},
		},
		{
			output: "/end:PT0H/height:100/period:300/start:-P1D/stat:Average/width:600",
			input: &MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Average",
				"height": int64(100),
				"width":  int64(600),
			},
		},
		{
			output: "/end:PT0H/height:400/period:300/start:-P1D/stat:Average/width:200",
			input: &MetricsRequest{
				"start":  "-P1D",
				"end":    "PT0H",
				"period": int64(300),
				"stat":   "Average",
				"height": int64(400),
				"width":  int64(200),
			},
		},
		{
			output: "/end:PT0H/height:2000/metrics:[AWS/ECS CPUUtilization ClusterName spinup-000393 ServiceName spinup-0010a3-testsvc]/period:300/start:-P1D/stat:Average/width:2000",
			input: &MetricsRequest{
				"start":   "-P1D",
				"end":     "PT0H",
				"period":  int64(300),
				"stat":    "Average",
				"height":  int64(2000),
				"width":   int64(2000),
				"metrics": []string{"AWS/ECS", "CPUUtilization", "ClusterName", "spinup-000393", "ServiceName", "spinup-0010a3-testsvc"},
			},
		},
	}

	for _, test := range tests {
		if out := test.input.String(); out != test.output {
			t.Errorf("expected '%s', got '%s'", test.output, out)
		}
	}
}
