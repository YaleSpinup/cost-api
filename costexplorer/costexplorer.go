package costexplorer

import (
	"github.com/YaleSpinup/cost-api/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/costexplorer/costexploreriface"
	log "github.com/sirupsen/logrus"
)

// CostExplorer is a wrapper around the aws costexplorer service with some default config info
type CostExplorer struct {
	Service         costexploreriface.CostExplorerAPI
}

// NewSession creates a new costexplorer session
func NewSession(account common.Account) CostExplorer {
	c := CostExplorer{}
	log.Infof("creating new aws session for costexplorer with key id %s in region %s", account.Akid, account.Region)
	sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(account.Akid, account.Secret, ""),
			Region:      aws.String(account.Region),
	}))
	c.Service = costexplorer.New(sess)
	return c
}