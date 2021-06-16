package computeoptimizer

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/computeoptimizer"
	"github.com/aws/aws-sdk-go/service/computeoptimizer/computeoptimizeriface"
)

// mockComputeOptimizerClient is a fake computeoptimizer client
type mockComputeOptimizerClient struct {
	computeoptimizeriface.ComputeOptimizerAPI
	t   *testing.T
	err error
}

func newMockComputeOptimizerClient(t *testing.T, err error) computeoptimizeriface.ComputeOptimizerAPI {
	return &mockComputeOptimizerClient{
		t:   t,
		err: err,
	}
}

func TestNewSession(t *testing.T) {
	client := New()
	to := reflect.TypeOf(client).String()
	if to != "*computeoptimizer.ComputeOptimizer" {
		t.Errorf("expected type to be '*computeoptimizer.ComputeOptimizer', got %s", to)
	}
}

var recommendation = `
[
    {
        "AccountId": "1234567890",
        "CurrentInstanceType": "m4.xlarge",
        "Finding": "UNDER_PROVISIONED",
        "FindingReasonCodes": [
            "CPUOverprovisioned",
            "EBSIOPSOverprovisioned",
            "EBSThroughputUnderprovisioned"
        ],
        "InstanceArn": "arn:aws:ec2:us-east-1:1234567890:instance/i-0987654321",
        "InstanceName": "knowledge.bobs.edu",
        "LastRefreshTimestamp": "2021-06-16T18:53:25.669Z",
        "LookBackPeriodInDays": 14,
        "RecommendationOptions": [
            {
                "InstanceType": "t3.xlarge",
                "PerformanceRisk": 3,
                "PlatformDifferences": [
                    "NetworkInterface",
                    "Hypervisor",
                    "StorageInterface"
                ],
                "ProjectedUtilizationMetrics": [
                    {
                        "Name": "CPU",
                        "Statistic": "MAXIMUM",
                        "Value": 44.64812085482684
                    }
                ],
                "Rank": 1
            },
            {
                "InstanceType": "m5.xlarge",
                "PerformanceRisk": 1,
                "PlatformDifferences": [
                    "NetworkInterface",
                    "Hypervisor",
                    "StorageInterface"
                ],
                "ProjectedUtilizationMetrics": [
                    {
                        "Name": "CPU",
                        "Statistic": "MAXIMUM",
                        "Value": 44.64812085482684
                    }
                ],
                "Rank": 2
            },
            {
                "InstanceType": "m4.xlarge",
                "PerformanceRisk": 1,
                "PlatformDifferences": [],
                "ProjectedUtilizationMetrics": [
                    {
                        "Name": "CPU",
                        "Statistic": "MAXIMUM",
                        "Value": 55.5084745762712
                    }
                ],
                "Rank": 3
            }
        ],
        "RecommendationSources": [
            {
                "RecommendationSourceArn": "arn:aws:ec2:us-east-1:1234567890:instance/i-0987654321",
                "RecommendationSourceType": "Ec2Instance"
            }
        ],
        "UtilizationMetrics": [
            {
                "Name": "CPU",
                "Statistic": "MAXIMUM",
                "Value": 55.5084745762712
            },
            {
                "Name": "EBS_READ_OPS_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 753.56
            },
            {
                "Name": "EBS_WRITE_OPS_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 917.98
            },
            {
                "Name": "EBS_READ_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 89428489.58333333
            },
            {
                "Name": "EBS_WRITE_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 83148230.79666667
            },
            {
                "Name": "NETWORK_IN_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 414224.09777777776
            },
            {
                "Name": "NETWORK_OUT_BYTES_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 29260.118055555555
            },
            {
                "Name": "NETWORK_PACKETS_IN_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 177.9728888888889
            },
            {
                "Name": "NETWORK_PACKETS_OUT_PER_SECOND",
                "Statistic": "MAXIMUM",
                "Value": 35.34688888888889
            }
        ]
    }
]
`

func unmarshallreq(req string) []*computeoptimizer.InstanceRecommendation {
	reqs := []*computeoptimizer.InstanceRecommendation{}
	json.Unmarshal([]byte(recommendation), &reqs)
	return reqs
}

func (m mockComputeOptimizerClient) GetEC2InstanceRecommendationsWithContext(ctx context.Context, input *computeoptimizer.GetEC2InstanceRecommendationsInput, opts ...request.Option) (*computeoptimizer.GetEC2InstanceRecommendationsOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	reqs := []*computeoptimizer.InstanceRecommendation{}
	if aws.StringValue(input.InstanceArns[0]) == "arn:aws:ec2:us-east-1:1234567890:instance/i-0987654321" {
		reqs = unmarshallreq(recommendation)
	}

	return &computeoptimizer.GetEC2InstanceRecommendationsOutput{
		Errors:                  []*computeoptimizer.GetRecommendationError{},
		InstanceRecommendations: reqs,
	}, nil
}

func TestComputeOptimizer_GetEc2InstanceRecommendations(t *testing.T) {
	type fields struct {
		session *session.Session
		Service computeoptimizeriface.ComputeOptimizerAPI
	}
	type args struct {
		ctx context.Context
		arn string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*computeoptimizer.InstanceRecommendation
		wantErr bool
	}{
		{
			name: "empty input",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "recommendations",
			fields: fields{
				Service: newMockComputeOptimizerClient(t, nil),
			},
			args: args{
				ctx: context.TODO(),
				arn: "arn:aws:ec2:us-east-1:1234567890:instance/i-0987654321",
			},
			want: unmarshallreq(recommendation),
		},
		{
			name: "empty recommendations",
			fields: fields{
				Service: newMockComputeOptimizerClient(t, nil),
			},
			args: args{
				ctx: context.TODO(),
				arn: "foobar",
			},
			want: []*computeoptimizer.InstanceRecommendation{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComputeOptimizer{
				session: tt.fields.session,
				Service: tt.fields.Service,
			}
			got, err := c.GetEc2InstanceRecommendations(tt.args.ctx, tt.args.arn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComputeOptimizer.GetEc2InstanceRecommendations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeOptimizer.GetEc2InstanceRecommendations() = %v, want %v", got, tt.want)
			}
		})
	}
}
