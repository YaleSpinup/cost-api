package api

import (
	"context"
	"fmt"
	"time"

	ce "github.com/YaleSpinup/cost-api/costexplorer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type costAndUsageReq struct {
	account, spaceID, start, end, groupBy string
}

func (o *costExplorerOrchestrator) getCostAndUsageForSpace(ctx context.Context, req *costAndUsageReq) ([]*costexplorer.ResultByTime, bool, time.Duration, error) {
	start, end, err := parseTime(req.start, req.end)
	if err != nil {
		return nil, false, 0, nil
	}

	input := costexplorer.GetCostAndUsageInput{
		Filter:      ce.And(inSpace(req.spaceID), inOrg(o.server.org), notTryIT()),
		Granularity: aws.String("MONTHLY"),
		Metrics: []*string{
			aws.String("BLENDED_COST"),
			aws.String("UNBLENDED_COST"),
			aws.String("USAGE_QUANTITY"),
		},
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(start),
			End:   aws.String(end),
		},
	}

	switch {
	case req.groupBy == "RESOURCE_NAME":
		input.GroupBy = []*costexplorer.GroupDefinition{
			{
				Key:  aws.String("Name"),
				Type: aws.String("TAG"),
			},
		}
	case req.groupBy != "":
		input.GroupBy = []*costexplorer.GroupDefinition{
			{
				Key:  aws.String(req.groupBy),
				Type: aws.String("DIMENSION"),
			},
		}
	}

	// create a cacheKey more unique than spaceID for managing cache objects.
	// Since we will accept date-range cost exploring and grouping, concatenate
	// the spaceID, the start time, end time and group by so we can cache each
	// time-based result
	cacheKey := fmt.Sprintf("%s_%s_%s_%s_%s", req.account, req.spaceID, req.start, req.end, req.groupBy)

	log.Debugf("cacheKey: %s", cacheKey)

	// the object is not found in the cache, call AWS cost-explorer and set cache
	c, expire, ok := o.server.resultCache.GetWithExpiration(cacheKey)
	if !ok || c == nil {
		log.Debugf("cache empty for org, and space-cacheKey: %s, %s, calling cost-explorer", o.server.org, cacheKey)

		// call cost-explorer
		out, err := o.client.GetCostAndUsage(ctx, &input)
		if err != nil {
			return nil, false, 0, err
		}

		// cache results
		o.server.resultCache.SetDefault(cacheKey, out)

		return out, false, 0, nil
	}

	// cached object was found
	out, ok := c.([]*costexplorer.ResultByTime)
	if !ok {
		return nil, false, 0, errors.New("value in cache is not a []*costexplorer.ResultByTime!")
	}

	log.Debugf("found cached object: %s", out)

	return out, true, time.Until(expire), nil
}

// parseTime returns time range from beginning of month to day-of-month now if the
// passed values are empty otherwise, it parses the string and returns the value (or an error)
func parseTime(start, end string) (string, string, error) {
	// if it's the first day of the month, get today's usage thus far
	// TODO: :confused:
	y, m, d := time.Now().Date()
	if d == 1 {
		d = 3
	}

	if start == "" {
		start = fmt.Sprintf("%d-%02d-01", y, m)
	}

	startStamp, err := time.Parse("2006-01-02", start)
	if err != nil {
		return "", "", err
	}

	if end == "" {
		end = fmt.Sprintf("%d-%02d-%02d", y, m, d)
	}

	endStamp, err := time.Parse("2006-01-02", end)
	if err != nil {
		return "", "", err
	}

	if !endStamp.After(startStamp) {
		return "", "", fmt.Errorf("end time should be after start time")
	}

	// convert time.Time to a string
	return startStamp.Format("2006-01-02"), endStamp.Format("2006-01-02"), nil
}

// inSpace returns the cost explorer expression to filter on spaceid
func inSpace(spaceID string) *costexplorer.Expression {
	return ce.Tag("spinup:spaceid", []string{spaceID})
}

// ofName returns the cost explorer expression to filter on name
func ofName(name string) *costexplorer.Expression {
	return ce.Tag("Name", []string{name})
}

// inOrg returns the cost explorer expression to filter on org
func inOrg(org string) *costexplorer.Expression {
	yaleTag := ce.Tag("yale:org", []string{org})
	spinupTag := ce.Tag("spinup:org", []string{org})
	return ce.Or(yaleTag, spinupTag)
}

// notTryIT returns the cost explorer expression to filter out tryits
func notTryIT() *costexplorer.Expression {
	yaleTag := ce.Tag("yale:subsidized", []string{"true"})
	spinupTag := ce.Tag("spinup:subsidized", []string{"true"})
	return ce.Not(ce.Or(yaleTag, spinupTag))
}
