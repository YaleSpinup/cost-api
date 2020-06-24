package cloudwatch

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	log "github.com/sirupsen/logrus"
)

type Metric []string
type MetricsRequest map[string]interface{}

// GetMetricWidget gets a metric widget image for an instance id
func (c *Cloudwatch) GetMetricWidget(ctx context.Context, req MetricsRequest) ([]byte, error) {
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

func (m *MetricsRequest) String() string {
	var s string
	if m == nil {
		return s
	}

	req := map[string]interface{}(*m)

	keys := make([]string, 0, len(req))
	for k := range req {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		s += fmt.Sprintf("/%s:%v", k, req[k])
	}

	return s
}
