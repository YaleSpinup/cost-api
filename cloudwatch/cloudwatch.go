package cloudwatch

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	log "github.com/sirupsen/logrus"
)

// Cloudwatch is a wrapper around the aws cloudwatch service with some default config info
type Cloudwatch struct {
	session *session.Session
	Service cloudwatchiface.CloudWatchAPI
}

type CloudwatchOption func(*Cloudwatch)

func New(opts ...CloudwatchOption) *Cloudwatch {
	client := Cloudwatch{}

	for _, opt := range opts {
		opt(&client)
	}

	if client.session != nil {
		client.Service = cloudwatch.New(client.session)
	}

	return &client
}

func WithSession(sess *session.Session) CloudwatchOption {
	return func(client *Cloudwatch) {
		log.Debug("using aws session")
		client.session = sess
	}
}

func WithCredentials(key, secret, token, region string) CloudwatchOption {
	return func(client *Cloudwatch) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		client.session = sess
	}
}
