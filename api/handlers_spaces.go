package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

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
	startTime := vars["StartTime"]
	endTime := vars["EndTime"]
	spaceID := vars["space"]

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
	var timeValidity bool
	var start string
	var end string
	// Did we get cost-explorer start and end times on the API?
	// set defaults, else verify times given on API
	if endTime == "" || startTime == "" {
		log.Debug("no start or end time given on API input, assigning defaults")
		// if it's the first day of the month, get today's usage thus far
		y, m, d := time.Now().Date()
		if d == 1 {
			d = 2
		}

		start = fmt.Sprintf("%d-%02d-01", y, m)
		end = fmt.Sprintf("%d-%02d-%02d", y, m, d)
		timeValidity = true
	} else {
		// debugging verbosity
		log.Debugf("startTime: %s\nendTime: %s\n", startTime, endTime)

		// sTmp and eTmp temporary vars to hold time.Time, then convert
		// both to strings
		sTmp, err := time.Parse("2006-01-02", startTime)
		if err != nil {
			msg := fmt.Sprintf("error parsing StartTime from input: %s\n", err)
			handleError(w, apierror.New(apierror.ErrBadRequest, msg, nil))
			timeValidity = false
			return
		}
		// convert time.Time to a string
		start = fmt.Sprint(sTmp.Format("2006-01-02"))

		eTmp, err := time.Parse("2006-01-02", endTime)
		if err != nil {
			msg := fmt.Sprintf("error parsing EndTime from input: %s\n", err)
			handleError(w, apierror.New(apierror.ErrBadRequest, msg, nil))
			timeValidity = false
			return
		}
		// convert time.Time to a string
		end = fmt.Sprint(eTmp.Format("2006-01-02"))

		// if time on the API input is already borked, don't continue
		// end time is greater than start time, logically
		timeValidity = eTmp.After(sTmp)
		log.Debugf("timeValidity value: %v", timeValidity)
		log.Debugf("startTime value: %s", startTime)
		log.Debugf("endTime value: %s", endTime)
		log.Debugf("start_inside_if value: %s", start)
		log.Debugf("end_inside_if value: %s", end)
		if !timeValidity {
			msg := fmt.Sprint("endTime should be greater that startTime\n")
			handleError(w, apierror.New(apierror.ErrBadRequest, msg, nil))
		}
	}

	log.Debugf("start_outside: %s\n", start)
	log.Debugf("end_outside: %s\n", end)

	input := costexplorer.GetCostAndUsageInput{
		Filter: &costexplorer.Expression{
			And: []*costexplorer.Expression{
				&costexplorer.Expression{
					Tags: &costexplorer.TagValues{
						Key: aws.String("spinup:spaceid"),
						Values: []*string{
							aws.String(spaceID),
						},
					},
				},
				&costexplorer.Expression{
					Or: []*costexplorer.Expression{
						&costexplorer.Expression{
							Tags: &costexplorer.TagValues{
								Key: aws.String("yale:org"),
								Values: []*string{
									aws.String(Org),
								},
							},
						},
						&costexplorer.Expression{
							Tags: &costexplorer.TagValues{
								Key: aws.String("spinup:org"),
								Values: []*string{
									aws.String(Org),
								},
							},
						},
					},
				},
				&costexplorer.Expression{
					Not: &costexplorer.Expression{
						Or: []*costexplorer.Expression{
							&costexplorer.Expression{
								Tags: &costexplorer.TagValues{
									Key: aws.String("yale:subsidized"),
									Values: []*string{
										aws.String("true"),
									},
								},
							},
							&costexplorer.Expression{
								Tags: &costexplorer.TagValues{
									Key: aws.String("spinup:subsidized"),
									Values: []*string{
										aws.String("true"),
									},
								},
							},
						},
					},
				},
			},
		},
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

	// if time did not parse correctly, don't call AWS cost explorer
	if timeValidity == true {
		// create a cacheKey more unique than spaceID for managing cache objects.
		// Since we will accept date-range cost exploring, concatenate the spaceID
		// and the start and end time so we can cache each time-based result
		var cacheKey string
		cacheKey = fmt.Sprintf("%s_%s_%s", spaceID, startTime, endTime)
		log.Debugf("cacheKey: %s", cacheKey)

		// the object is not found in the cache, call AWS cost-explorer and set cache
		var out []*costexplorer.ResultByTime
		c, expire, ok := resultCache.GetWithExpiration(cacheKey)
		if !ok || c == nil {
			log.Debugf("cache empty for org, and space-cacheKey: %s, %s, calling cost-explorer", Org, cacheKey)
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

}
