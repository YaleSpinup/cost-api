package resourcegroupstaggingapi

import (
	"context"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"
	log "github.com/sirupsen/logrus"
)

// ResourceGroupsTaggingAPI is a wrapper around the aws resourcegroupstaggingapi service with some default config info
type ResourceGroupsTaggingAPI struct {
	session *session.Session
	Service resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
}

type ResourceGroupsTaggingAPIOption func(*ResourceGroupsTaggingAPI)

func New(opts ...ResourceGroupsTaggingAPIOption) *ResourceGroupsTaggingAPI {
	client := ResourceGroupsTaggingAPI{}

	for _, opt := range opts {
		opt(&client)
	}

	if client.session != nil {
		client.Service = resourcegroupstaggingapi.New(client.session)
	}

	return &client
}

func WithSession(sess *session.Session) ResourceGroupsTaggingAPIOption {
	return func(client *ResourceGroupsTaggingAPI) {
		log.Debug("using aws session")
		client.session = sess
	}
}

func WithCredentials(key, secret, token, region string) ResourceGroupsTaggingAPIOption {
	return func(client *ResourceGroupsTaggingAPI) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		client.session = sess
	}
}

func (r *ResourceGroupsTaggingAPI) ListResourcesWithTags(ctx context.Context, input *resourcegroupstaggingapi.GetResourcesInput) (*resourcegroupstaggingapi.GetResourcesOutput, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Info("listing tagged resources")

	out, err := r.Service.GetResourcesWithContext(ctx, input)
	if err != nil {
		return nil, ErrCode("listing resource with tags", err)
	}

	log.Debugf("got output from get resources: %+v", out)

	return out, nil
}
