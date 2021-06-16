package api

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/computeoptimizer"
)

func (o *optimizerOrchestrator) GetInstanceRecommendations(ctx context.Context, account, id string) ([]*computeoptimizer.InstanceRecommendation, error) {
	a := arn.ARN{
		Partition: "aws",
		Region:    "us-east-1",
		AccountID: account,
		Service:   "ec2",
		Resource:  fmt.Sprintf("instance/%s", id),
	}

	out, err := o.client.GetEc2InstanceRecommendations(ctx, a.String())
	if err != nil {
		return nil, err
	}

	return out, nil
}
