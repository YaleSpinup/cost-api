package cloudwatch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	log "github.com/sirupsen/logrus"
)

// GetMetricWidget gets a metric widget image for an instance id
func (c *Cloudwatch) GetMetricWidget(ctx context.Context, metric, id string) ([]byte, error) {
	if metric == "" || id == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting metric widget for metic '%s' of instance '%s'", metric, id)

	// default to last day of metrics for one instance id
	req := struct {
		Metrics [][]string `json:"metrics"`
		Period  int64      `json:"period"`
		Start   string     `json:"start"`
		End     string     `json:"end"`
	}{
		Metrics: [][]string{
			[]string{
				"AWS/EC2",
				metric,
				"InstanceId",
				id,
			},
		},
		Period: 300,
		Start:  "-P1D",
		End:    "PT0H",
	}

	j, err := json.Marshal(req)
	if err != nil {
		msg := fmt.Sprintf("failed to get build widget request input for metric '%s' and instance id '%s': %s", metric, id, err)
		return nil, ErrCode(msg, err)
	}

	log.Debugf("getting metric widget with input request %s", string(j))

	out, err := c.Service.GetMetricWidgetImageWithContext(ctx, &cloudwatch.GetMetricWidgetImageInput{MetricWidget: aws.String(string(j))})
	if err != nil {
		msg := fmt.Sprintf("failed to get metric widget image for metric '%s' and instance id '%s': %s", metric, id, err)
		return nil, ErrCode(msg, err)
	}

	return out.MetricWidgetImage, nil
}
