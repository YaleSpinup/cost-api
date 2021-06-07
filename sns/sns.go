package sns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	log "github.com/sirupsen/logrus"
)

type SNS struct {
	session *session.Session
	Service snsiface.SNSAPI
}

type SNSOption func(*SNS)

func New(opts ...SNSOption) *SNS {
	client := SNS{}

	for _, opt := range opts {
		opt(&client)
	}

	if client.session != nil {
		client.Service = sns.New(client.session)
	}

	return &client
}

func WithSession(sess *session.Session) SNSOption {
	return func(client *SNS) {
		log.Debug("using aws session")
		client.session = sess
	}
}

func WithCredentials(key, secret, token, region string) SNSOption {
	return func(client *SNS) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		client.session = sess
	}
}
