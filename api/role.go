package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/YaleSpinup/aws-go/services/session"
	stsSvc "github.com/YaleSpinup/aws-go/services/sts"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/google/uuid"
	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

// assumeRole assumes the passed role arn.  if an externalId is set in the account to be accessed, it can be passed with the request.  inline
// policy can be passed to limit the access for the session.  policy Arns can also be passed to limit access for the session.
func (s *server) assumeRole(ctx context.Context, externalId, roleArn, inlinePolicy string, policyArns ...string) (*session.Session, error) {
	contextLogger := log.WithFields(log.Fields{
		"role": roleArn,
	})

	start := time.Now()
	defer func() {
		totalTime := time.Since(start)
		contextLogger.WithField("duration", totalTime).Info("assumeRole()")
	}()

	stsService := stsSvc.New(stsSvc.WithSession(s.session.Session))

	name := fmt.Sprintf("spinup-%s-cost-api-%s", s.org, uuid.New())

	contextLogger = contextLogger.WithField("session", name)

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

	cacheKey := fmt.Sprintf("spinup_%s_%s", s.org, roleArn)

	if externalId != "" {
		input.SetExternalId(externalId)
		cacheKey = cacheKey + "_" + externalId
	}

	if inlinePolicy != "" {
		input.SetPolicy(inlinePolicy)
		cacheKey = cacheKey + "_" + inlinePolicy
	}

	if policyArns != nil {
		arns := []*sts.PolicyDescriptorType{}
		for _, a := range policyArns {
			arns = append(arns, &sts.PolicyDescriptorType{
				Arn: aws.String(a),
			})
		}
		input.SetPolicyArns(arns)

		cacheKey = cacheKey + "_" + strings.Join(policyArns, "_")
	}

	contextLogger.Debugf("checking for item with cache key: '%s'", cacheKey)

	item, expire, found := s.sessionCache.GetWithExpiration(cacheKey)
	if found {
		if sess, ok := item.(*session.Session); ok {
			contextLogger.Infof("using cached session (expire: %s)", expire.String())
			return sess, nil
		}
	}

	contextLogger.Debugf("assuming role %s with input %+v", roleArn, input)

	out, err := stsService.AssumeRole(ctx, &input)
	if err != nil {
		log.Errorf("got: %s", err)
		return nil, err
	}

	akid := aws.StringValue(out.Credentials.AccessKeyId)

	contextLogger.Infof("got temporary creds %s, expiration: %s", akid, aws.TimeValue(out.Credentials.Expiration).String())

	sess := session.New(
		session.WithCredentials(
			akid,
			aws.StringValue(out.Credentials.SecretAccessKey),
			aws.StringValue(out.Credentials.SessionToken),
		),
		session.WithRegion("us-east-1"),
	)

	contextLogger.Debugf("caching session with cache key: '%s'", cacheKey)

	s.sessionCache.Set(cacheKey, &sess, cache.DefaultExpiration)

	return &sess, nil
}
