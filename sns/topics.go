package sns

import (
	"context"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	log "github.com/sirupsen/logrus"
)

// CreateTopic creates an SNS topic
func (s *SNS) CreateTopic(ctx context.Context, input *sns.CreateTopicInput) (*sns.CreateTopicOutput, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("creating SNS topic %s", aws.StringValue(input.Name))

	out, err := s.Service.CreateTopicWithContext(ctx, input)
	if err != nil {
		return nil, ErrCode("failed to create sns topic", err)
	}

	log.Debugf("got output creating topic: %+v", out)

	return out, nil
}

// DeleteTopic deletes an SNS topic
func (s *SNS) DeleteTopic(ctx context.Context, arn string) error {
	if arn == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("deleting sns topic %s", arn)

	if _, err := s.Service.DeleteTopicWithContext(ctx, &sns.DeleteTopicInput{
		TopicArn: aws.String(arn),
	}); err != nil {
		return ErrCode("failed to delete sns topic", err)
	}

	return nil
}

func (s *SNS) CreateSubscription(ctx context.Context, input *sns.SubscribeInput) (*sns.SubscribeOutput, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("subscribing to SNS topic %s", aws.StringValue(input.TopicArn))

	out, err := s.Service.SubscribeWithContext(ctx, input)
	if err != nil {
		return nil, ErrCode("failed to subscribe to sns topic", err)
	}

	log.Debugf("got output subscribing to topic: %+v", out)

	return out, nil
}
