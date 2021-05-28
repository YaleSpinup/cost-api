package api

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/aws-go/services/session"
	stsSvc "github.com/YaleSpinup/aws-go/services/sts"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// assumeRole assumes the passed role arn.  if an externalId is set in the account to be accessed, it can be passed with the request.  inline
// policy can be passed to limit the access for the session.  policy Arns can also be passed to limit access for the session.
func (s *server) assumeRole(ctx context.Context, externalId, roleArn, inlinePolicy string, policyArns ...string) (*session.Session, error) {
	log.Infof("server assuming role %s", roleArn)

	stsService := stsSvc.New(stsSvc.WithSession(s.session.Session))

	name := fmt.Sprintf("spinup-%s-ecr-api-%s", s.org, uuid.New())

	input := sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(900),
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(name),
		Tags: []*sts.Tag{
			{
				Key:   aws.String("spinup:org"),
				Value: aws.String(s.org),
			},
		},
	}

	if externalId != "" {
		input.SetExternalId(externalId)
	}

	if inlinePolicy != "" {
		input.SetPolicy(inlinePolicy)
	}

	if policyArns != nil {
		arns := []*sts.PolicyDescriptorType{}
		for _, a := range policyArns {
			arns = append(arns, &sts.PolicyDescriptorType{
				Arn: aws.String(a),
			})
		}
		input.SetPolicyArns(arns)
	}

	out, err := stsService.AssumeRole(ctx, &input)
	if err != nil {
		log.Errorf("got: %s", err)
		return nil, err
	}

	akid := aws.StringValue(out.Credentials.AccessKeyId)

	log.Infof("got temporary creds %s, expiration: %s", akid, aws.TimeValue(out.Credentials.Expiration).String())

	sess := session.New(
		session.WithCredentials(
			akid,
			aws.StringValue(out.Credentials.SecretAccessKey),
			aws.StringValue(out.Credentials.SessionToken),
		),
		session.WithRegion("us-east-1"),
	)

	return &sess, nil
}
