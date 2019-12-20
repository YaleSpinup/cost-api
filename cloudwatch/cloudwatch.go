package cloudwatch

import (
	"github.com/YaleSpinup/cost-api/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	log "github.com/sirupsen/logrus"
)

// Cloudwatch is a wrapper around the aws cloudwatch service with some default config info
type Cloudwatch struct {
	Service cloudwatchiface.CloudWatchAPI
}

// NewSession creates a new cloudwatch session
func NewSession(account common.Account) Cloudwatch {
	c := Cloudwatch{}
	log.Infof("creating new aws session for costexplorer with key id %s in region %s", account.Akid, account.Region)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(account.Akid, account.Secret, ""),
		Region:      aws.String(account.Region),
	}))
	c.Service = cloudwatch.New(sess)
	return c
}
