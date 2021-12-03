package costexplorer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/costexplorer/costexploreriface"
	log "github.com/sirupsen/logrus"
)

// CostExplorer is a wrapper around the aws costexplorer service with some default config info
type CostExplorer struct {
	session *session.Session
	Service costexploreriface.CostExplorerAPI
}

type CostExplorerOption func(*CostExplorer)

func New(opts ...CostExplorerOption) *CostExplorer {
	client := CostExplorer{}

	for _, opt := range opts {
		opt(&client)
	}

	if client.session != nil {
		client.Service = costexplorer.New(client.session)
	}

	return &client
}

func WithSession(sess *session.Session) CostExplorerOption {
	return func(client *CostExplorer) {
		log.Debug("using aws session")
		client.session = sess
	}
}

func WithCredentials(key, secret, token, region string) CostExplorerOption {
	return func(client *CostExplorer) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		client.session = sess
	}
}
