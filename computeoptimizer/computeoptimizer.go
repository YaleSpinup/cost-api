package computeoptimizer

import (
	"context"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/computeoptimizer"
	"github.com/aws/aws-sdk-go/service/computeoptimizer/computeoptimizeriface"
	log "github.com/sirupsen/logrus"
)

type ComputeOptimizer struct {
	session *session.Session
	Service computeoptimizeriface.ComputeOptimizerAPI
}

type ComputeOptimizerOption func(*ComputeOptimizer)

func New(opts ...ComputeOptimizerOption) *ComputeOptimizer {
	client := ComputeOptimizer{}

	for _, opt := range opts {
		opt(&client)
	}

	if client.session != nil {
		client.Service = computeoptimizer.New(client.session)
	}

	return &client
}

func WithSession(sess *session.Session) ComputeOptimizerOption {
	return func(client *ComputeOptimizer) {
		log.Debug("using aws session")
		client.session = sess
	}
}

func WithCredentials(key, secret, token, region string) ComputeOptimizerOption {
	return func(client *ComputeOptimizer) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		client.session = sess
	}
}

func (c *ComputeOptimizer) GetEc2InstanceRecommendations(ctx context.Context, arn string) ([]*computeoptimizer.InstanceRecommendation, error) {
	if arn == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting recommendations for ec2 instance %s", arn)

	out, err := c.Service.GetEC2InstanceRecommendationsWithContext(ctx, &computeoptimizer.GetEC2InstanceRecommendationsInput{
		InstanceArns: aws.StringSlice([]string{arn}),
	})
	if err != nil {
		return nil, ErrCode("failed to get instance recommendations", err)
	}

	log.Debugf("got output from ec2 instance recommendations for %s: %+v", arn, out)

	return out.InstanceRecommendations, nil
}
