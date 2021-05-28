package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	ce "github.com/YaleSpinup/cost-api/costexplorer"
	log "github.com/sirupsen/logrus"
)

// SpaceGetHandler gets the cost for a space, grouped by the service.  By default,
// it pulls data from the start of the month until now.
func (s *server) SpaceGetHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}

	// loop thru and log given API input vars in debug
	vars := mux.Vars(r)
	for k, v := range vars {
		log.Debugf("key=%v, value=%v", k, v)
	}

	// get vars from the API route
	account := vars["account"]
	startTime := vars["start"]
	endTime := vars["end"]
	spaceID := vars["space"]
	groupBy := vars["groupby"]

	ceService, ok := s.costExplorerServices[account]
	if !ok {
		msg := fmt.Sprintf("cost explorer service not found for account: %s", account)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}
	log.Debugf("found cost explorer service %+v", ceService)

	resultCache, ok := s.resultCache[account]
	if !ok {
		msg := fmt.Sprintf("result cache not found for account: %s", account)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}
	log.Debugf("found cost explorer result cache %+v", *resultCache)

	// vars for checking input times parse and are valid
	var start string
	var end string
	var err error

	// Did we get cost-explorer start and end times on the API?
	// set defaults, else verify times given on API
	if endTime == "" || startTime == "" {
		log.Debug("no start or end time given on API input, assigning defaults")
		start, end = getTimeDefault()
	} else {
		start, end, err = getTimeAPI(startTime, endTime)
		if err != nil {
			handleError(w, err)
			return
		}
	}

	input := costexplorer.GetCostAndUsageInput{
		Filter:      ce.And(inSpace(spaceID), inOrg(s.org), notTryIT()),
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

	if groupBy != "" {
		input.GroupBy = []*costexplorer.GroupDefinition{
			{
				Key:  aws.String(groupBy),
				Type: aws.String("DIMENSION"),
			},
		}
	}

	// create a cacheKey more unique than spaceID for managing cache objects.
	// Since we will accept date-range cost exploring and grouping, concatenate
	// the spaceID, the start time, end time and group by so we can cache each
	// time-based result
	cacheKey := fmt.Sprintf("%s_%s_%s_%s", spaceID, startTime, endTime, groupBy)
	log.Debugf("cacheKey: %s", cacheKey)

	// the object is not found in the cache, call AWS cost-explorer and set cache
	var out []*costexplorer.ResultByTime
	c, expire, ok := resultCache.GetWithExpiration(cacheKey)
	if !ok || c == nil {
		log.Debugf("cache empty for org, and space-cacheKey: %s, %s, calling cost-explorer", s.org, cacheKey)
		// call cost-explorer
		var err error
		out, err = ceService.GetCostAndUsage(r.Context(), &input)
		if err != nil {
			msg := fmt.Sprintf("failed to get costs for space %s: %s", cacheKey, err.Error())
			handleError(w, errors.Wrap(err, msg))
			return
		}

		// cache results
		resultCache.SetDefault(cacheKey, out)
	} else {
		// cached object was found
		out = c.([]*costexplorer.ResultByTime)
		log.Debugf("found cached object: %s", out)
		w.Header().Set("X-Cache-Hit", "true")
		w.Header().Set("X-Cache-Expire", fmt.Sprintf("%0.fs", time.Until(expire).Seconds()))
	}

	j, err := json.Marshal(out)
	if err != nil {
		log.Errorf("cannot marshal response (%v) into JSON: %s", out, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

// SpaceResourceGetHandler gets the cost for a resource within a space
func (s *server) SpaceResourceGetHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}

	log.Debugf("ResourceGetHandler Called")

	// loop thru and log given API input vars in debug
	vars := mux.Vars(r)
	for k, v := range vars {
		log.Debugf("key=%v, value=%v", k, v)
	}

	// get vars from the API route
	account := vars["account"]
	startTime := vars["start"]
	endTime := vars["end"]
	spaceID := vars["space"]
	name := vars["resourcename"]

	ceService, ok := s.costExplorerServices[account]
	if !ok {
		msg := fmt.Sprintf("cost explorer service not found for account: %s", account)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}
	log.Debugf("found cost explorer service %+v", ceService)

	resultCache, ok := s.resultCache[account]
	if !ok {
		msg := fmt.Sprintf("result cache not found for account: %s", account)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}
	log.Debugf("found cost explorer result cache %+v", *resultCache)

	// vars for checking input times parse and are valid
	var start string
	var end string
	var err error

	// Did we get cost-explorer start and end times on the API?
	// set defaults, else verify times given on API
	if endTime == "" || startTime == "" {
		log.Debug("no start or end time given on API input, assigning defaults")
		start, end = getTimeDefault()
	} else {
		start, end, err = getTimeAPI(startTime, endTime)
		if err != nil {
			handleError(w, err)
			return
		}
	}

	input := costexplorer.GetCostAndUsageInput{
		Filter:      ce.And(ofName(name), inSpace(spaceID), inOrg(s.org), notTryIT()),
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

	// create a cacheKey more unique than spaceID for managing cache objects.
	// Since we will accept date-range cost exploring, concatenate the spaceID
	// and the start and end time so we can cache each time-based result
	cacheKey := fmt.Sprintf("%s_%s_%s_%s", spaceID, name, startTime, endTime)
	log.Debugf("cacheKey: %s", cacheKey)

	// the object is not found in the cache, call AWS cost-explorer and set cache
	var out []*costexplorer.ResultByTime
	c, expire, ok := resultCache.GetWithExpiration(cacheKey)
	if !ok || c == nil {
		log.Debugf("cache empty for org, and space-cacheKey: %s, %s, calling cost-explorer", s.org, cacheKey)
		// call cost-explorer
		var err error
		out, err = ceService.GetCostAndUsage(r.Context(), &input)
		if err != nil {
			msg := fmt.Sprintf("failed to get costs for space %s: %s", cacheKey, err.Error())
			handleError(w, errors.Wrap(err, msg))
			return
		}

		// cache results
		resultCache.SetDefault(cacheKey, out)
	} else {
		// cached object was found
		out = c.([]*costexplorer.ResultByTime)
		log.Debugf("found cached object: %s", out)
		w.Header().Set("X-Cache-Hit", "true")
		w.Header().Set("X-Cache-Expire", fmt.Sprintf("%0.fs", time.Until(expire).Seconds()))
	}

	j, err := json.Marshal(out)
	if err != nil {
		log.Errorf("cannot marshal response (%v) into JSON: %s", out, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

// getTimeDefault returns time range from beginning of month to day-of-month now
func getTimeDefault() (string, string) {
	// if it's the first day of the month, get today's usage thus far
	y, m, d := time.Now().Date()
	if d == 1 {
		d = 3
	}
	return fmt.Sprintf("%d-%02d-01", y, m), fmt.Sprintf("%d-%02d-%02d", y, m, d)
}

// getTimeAPI returns time parsed from API input
func getTimeAPI(startTime, endTime string) (string, string, error) {
	log.Debugf("startTime: %s, endTime: %s ", startTime, endTime)

	// sTmp and eTmp temporary vars to hold time.Time objects
	sTmp, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return "", "", errors.Wrapf(err, "error parsing StartTime from input")
	}

	eTmp, err := time.Parse("2006-01-02", endTime)
	if err != nil {
		return "", "", errors.Wrapf(err, "error parsing EndTime from input")
	}

	// if time on the API input is already borked, don't continue
	// end time is greater than start time, logically
	timeValidity := eTmp.After(sTmp)
	if !timeValidity {
		return "", "", errors.Errorf("endTime should be greater than startTime")
	}

	// convert time.Time to a string
	return sTmp.Format("2006-01-02"), eTmp.Format("2006-01-02"), nil
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
