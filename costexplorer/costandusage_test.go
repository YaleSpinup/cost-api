package costexplorer

import (
	"context"
	"reflect"
	"testing"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

var testResult1 = &costexplorer.ResultByTime{
	Estimated: aws.Bool(true),
	Groups: []*costexplorer.Group{
		{
			Keys: []*string{
				aws.String("Application$awesomeSauce"),
			},
			Metrics: map[string]*costexplorer.MetricValue{
				"BlendedCost": {
					Amount: aws.String("69.3586765877"),
					Unit:   aws.String("USD"),
				},
				"UsageQuanitity": {
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

func TestAnd(t *testing.T) {
	expected := &costexplorer.Expression{
		And: []*costexplorer.Expression{
			{
				Tags: &costexplorer.TagValues{
					Key: aws.String("mygroup"),
					Values: []*string{
						aws.String("superwings"),
					},
				},
			},
			{
				Tags: &costexplorer.TagValues{
					Key: aws.String("name"),
					Values: []*string{
						aws.String("jett"),
						aws.String("donnie"),
						aws.String("dizzy"),
						aws.String("jerome"),
					},
				},
			},
		},
	}
	out := And(Tag("mygroup", []string{"superwings"}), Tag("name", []string{"jett", "donnie", "dizzy", "jerome"}))
	if !awsutil.DeepEqual(expected, out) {
		t.Errorf("expected expression %s, got %s", awsutil.Prettify(expected), awsutil.Prettify(out))
	}
}

func TestOr(t *testing.T) {
	expected := &costexplorer.Expression{
		Or: []*costexplorer.Expression{
			{
				Tags: &costexplorer.TagValues{
					Key: aws.String("human"),
					Values: []*string{
						aws.String("jimbo"),
					},
				},
			},
			{
				Tags: &costexplorer.TagValues{
					Key: aws.String("human"),
					Values: []*string{
						aws.String("sky"),
					},
				},
			},
		},
	}
	out := Or(Tag("human", []string{"jimbo"}), Tag("human", []string{"sky"}))
	if !awsutil.DeepEqual(expected, out) {
		t.Errorf("expected expression %s, got %s", awsutil.Prettify(expected), awsutil.Prettify(out))
	}
}

func TestNot(t *testing.T) {
	expected := &costexplorer.Expression{
		Not: &costexplorer.Expression{
			Tags: &costexplorer.TagValues{
				Key: aws.String("delivery_time"),
				Values: []*string{
					aws.String("late"),
				},
			},
		},
	}
	out := Not(Tag("delivery_time", []string{"late"}))
	if !awsutil.DeepEqual(expected, out) {
		t.Errorf("expected expression %s, got %s", awsutil.Prettify(expected), awsutil.Prettify(out))
	}
}

func TestFilterChain(t *testing.T) {
	expected := &costexplorer.Expression{
		And: []*costexplorer.Expression{
			{
				Tags: &costexplorer.TagValues{
					Key: aws.String("mygroup"),
					Values: []*string{
						aws.String("superwings"),
					},
				},
			},
			{
				Or: []*costexplorer.Expression{
					{
						Tags: &costexplorer.TagValues{
							Key: aws.String("team"),
							Values: []*string{
								aws.String("rescue_riders"),
							},
						},
					},
					{
						Tags: &costexplorer.TagValues{
							Key: aws.String("team"),
							Values: []*string{
								aws.String("galaxy_wings"),
							},
						},
					},
				},
			},
			{
				Not: &costexplorer.Expression{
					Or: []*costexplorer.Expression{
						{
							Tags: &costexplorer.TagValues{
								Key: aws.String("name"),
								Values: []*string{
									aws.String("crystal"),
								},
							},
						},
						{
							Tags: &costexplorer.TagValues{
								Key: aws.String("name"),
								Values: []*string{
									aws.String("jerry"),
								},
							},
						},
					},
				},
			},
		},
	}
	notNew := Not(Or(Tag("name", []string{"crystal"}), Tag("name", []string{"jerry"})))
	inTeam := Or(Tag("team", []string{"rescue_riders"}), Tag("team", []string{"galaxy_wings"}))
	out := And(Tag("mygroup", []string{"superwings"}), inTeam, notNew)
	if !awsutil.DeepEqual(expected, out) {
		t.Errorf("expected expression %s, got %s", awsutil.Prettify(expected), awsutil.Prettify(out))
	}
}
