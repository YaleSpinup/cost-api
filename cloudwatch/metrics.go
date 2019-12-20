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

// GetMetricWidget gets a metric widget image for an instance id
func (c *Cloudwatch) GetMetricWidget(ctx context.Context, metrics []Metric, period int64, start, end string) ([]byte, error) {
	if len(metrics) == 0 || period == 0 || start == "" || end == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting metric widget for metics '%+v' with period %d, start %s and end %s", metrics, period, start, end)

	// default to last day of metrics for one instance id
	req := struct {
		Metrics []Metric `json:"metrics"`
		Period  int64    `json:"period"`
		Start   string   `json:"start"`
		End     string   `json:"end"`
	}{
		Metrics: metrics,
		Period:  period,
		Start:   start,
		End:     end,
	}

	j, err := json.Marshal(req)
	if err != nil {
		msg := fmt.Sprintf("failed to get build widget request input for metrics '%+v' with period %d, start %s and end %s: %s", metrics, period, start, end, err)
		return nil, ErrCode(msg, err)
	}

	log.Debugf("getting metric widget with input request %s", string(j))

	out, err := c.Service.GetMetricWidgetImageWithContext(ctx, &cloudwatch.GetMetricWidgetImageInput{MetricWidget: aws.String(string(j))})
	if err != nil {
		msg := fmt.Sprintf("failed to get metric widget image for metrics '%+v' with period %d, start %s and end %s: %s", metrics, period, start, end, err)
		return nil, ErrCode(msg, err)
	}

	return out.MetricWidgetImage, nil
}
