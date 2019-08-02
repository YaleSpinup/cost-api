package costexplorer

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	log "github.com/sirupsen/logrus"
)

// GetCostAndUsage gets cost and usage information from the cost explorer service
func (c *CostExplorer) GetCostAndUsage(ctx context.Context, input *costexplorer.GetCostAndUsageInput) ([]*costexplorer.ResultByTime, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting cost and usage with %+v", input)

	out, err := c.Service.GetCostAndUsageWithContext(ctx, input)
	if err != nil {
		msg := fmt.Sprintf("failed to get cost and usage report %+v", *input)
		return nil, ErrCode(msg, err)
	}

	log.Debugf("got cost and usage: %+v", out)

	return out.ResultsByTime, nil
}
