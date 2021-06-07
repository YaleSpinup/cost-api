package api

import (
	"encoding/json"

	"github.com/YaleSpinup/aws-go/services/iam"

	log "github.com/sirupsen/logrus"
)

// orgTagAccessPolicy generates the org tag conditional policy to be passed inline when assuming a role
func orgTagAccessPolicy(org string) (string, error) {
	log.Debugf("generating org policy document")

	policy := iam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iam.StatementEntry{
			{
				Effect:   "Allow",
				Action:   []string{"*"},
				Resource: []string{"*"},
				Condition: iam.Condition{
					"StringEquals": iam.ConditionStatement{
						"aws:ResourceTag/spinup:org": []string{org},
					},
				},
			},
		},
	}

	j, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func budgetReadWritePolicy() (string, error) {
	log.Debugf("generating org policy document")

	policy := iam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iam.StatementEntry{
			{
				Effect: "Allow",
				Action: []string{
					"budgets:ViewBudget",
					"budgets:ModifyBudget",
					"SNS:CreateTopic",
					"SNS:DeleteTopic",
					"SNS:Subscribe",
				},
				Resource: []string{"*"},
			},
		},
	}

	j, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}

	return string(j), nil

}

func defaultBudgetTopicPolicy(arn string) (string, error) {
	policy := iam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iam.StatementEntry{
			{
				Sid:    "AWSBudgetsSNSPublishingPermissions",
				Effect: "Allow",
				Principal: iam.Principal{
					"Service": iam.Value{
						"budgets.amazonaws.com",
					},
				},
				Action:   []string{"SNS:Publish"},
				Resource: []string{arn},
			},
		},
	}

	j, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}

	return string(j), nil
}
