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

type Metric []string
type MetricsRequest map[string]interface{}

// GetMetricWidget gets a metric widget image for an instance id
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/CloudWatch-Metric-Widget-Structure.html
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/viewing_metrics_with_cloudwatch.html
// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/cloudwatch-metrics.html
//
// Example metrics request
// {
//   "metrics": [
//     [ "AWS/ECS", "CPUUtilization", "ClusterName", "spinup-000393", "ServiceName", "spinup-0010a3-testsvc" ]
//   ],
//   "stat": "Average"
//   "period": 300,
//   "start": "-P1D",
//   "end": "PT0H"
// }
func (c *Cloudwatch) GetInstanceMetricWidget(ctx context.Context, req MetricsRequest) ([]byte, error) {
	if req == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting metric widget with input request %+v", req)

	j, err := json.Marshal(req)
	if err != nil {
		msg := fmt.Sprintf("failed to build widget request input for metrics from '%+v': %s", req, err)
		return nil, ErrCode(msg, err)
	}

	log.Debugf("getting metric widget with input json %s", string(j))

	out, err := c.Service.GetMetricWidgetImageWithContext(ctx, &cloudwatch.GetMetricWidgetImageInput{MetricWidget: aws.String(string(j))})
	if err != nil {
		msg := fmt.Sprintf("failed to get metric widget from request '%+s': %s", string(j), err)
		return nil, ErrCode(msg, err)
	}

	return out.MetricWidgetImage, nil
}
