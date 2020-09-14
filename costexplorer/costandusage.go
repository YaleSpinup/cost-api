package costexplorer

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	log "github.com/sirupsen/logrus"
)

// GetCostAndUsage gets cost and usage information from the cost explorer service
func (c *CostExplorer) GetCostAndUsage(ctx context.Context, input *costexplorer.GetCostAndUsageInput) ([]*costexplorer.ResultByTime, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting cost and usage with %+v", input)

	out, err := c.Service.GetCostAndUsageWithContext(ctx, input)
	if err != nil {
		msg := fmt.Sprintf("failed to get cost and usage report %+v", *input)
		return nil, ErrCode(msg, err)
	}

	log.Debugf("got cost and usage: %+v", out)

	return out.ResultsByTime, nil
}

// And returns the expressions wrapped in And
func And(exp ...*costexplorer.Expression) *costexplorer.Expression {
	expressions := []*costexplorer.Expression{}
	return &costexplorer.Expression{
		And: append(expressions, exp...),
	}
}

// Or returns the expressions wrapped in Or
func Or(exp ...*costexplorer.Expression) *costexplorer.Expression {
	expressions := []*costexplorer.Expression{}
	return &costexplorer.Expression{
		Or: append(expressions, exp...),
	}
}

// Not returns the negated expression
func Not(exp *costexplorer.Expression) *costexplorer.Expression {
	return &costexplorer.Expression{
		Not: exp,
	}
}

// Tag returns the cost explorer expression to filter on tag
func Tag(key string, values []string) *costexplorer.Expression {
	return &costexplorer.Expression{
		Tags: &costexplorer.TagValues{
			Key:    aws.String(key),
			Values: aws.StringSlice(values),
		},
	}
}
