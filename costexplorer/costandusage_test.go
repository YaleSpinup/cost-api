package costexplorer

import (
	"context"
	"reflect"
	"testing"

	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

var testResult1 = &costexplorer.ResultByTime{
	Estimated: aws.Bool(true),
	Groups: []*costexplorer.Group{
		&costexplorer.Group{
			Keys: []*string{
				aws.String("Application$awesomeSauce"),
			},
			Metrics: map[string]*costexplorer.MetricValue{
				"BlendedCost": &costexplorer.MetricValue{
					Amount: aws.String("69.3586765877"),
					Unit:   aws.String("USD"),
				},
				"UsageQuanitity": &costexplorer.MetricValue{
					Amount: aws.String("2088"),
					Unit:   aws.String("N/A"),
				},
			},
		},
	},
	TimePeriod: &costexplorer.DateInterval{
		End:   aws.String("2019-07-30"),
		Start: aws.String("2019-07-01"),
	},
	Total: map[string]*costexplorer.MetricValue{},
}

func (m *mockCostExplorerClient) GetCostAndUsageWithContext(ctx context.Context, input *costexplorer.GetCostAndUsageInput, opts ...request.Option) (*costexplorer.GetCostAndUsageOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &costexplorer.GetCostAndUsageOutput{
		ResultsByTime: []*costexplorer.ResultByTime{
			testResult1,
		},
	}, nil
}

func TestGetCostAndUsage(t *testing.T) {
	c := CostExplorer{
		Service: newmockCostExplorerClient(t, nil),
	}

	// test success
	expected := []*costexplorer.ResultByTime{testResult1}
	out, err := c.GetCostAndUsage(context.TODO(), &costexplorer.GetCostAndUsageInput{})
	if err != nil {
		t.Errorf("expected nil error, got: %s", err)
	}

	if !reflect.DeepEqual(out, expected) {
		t.Errorf("expected %+v, got %+v", expected, out)
	}

	// test nil input
	_, err = c.GetCostAndUsage(context.TODO(), nil)
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrBadRequest {
			t.Errorf("expected error code %s, got: %s", apierror.ErrBadRequest, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}
}
