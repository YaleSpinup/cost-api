package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/YaleSpinup/cost-api/cloudwatch"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MetricsGetImageHandler gets metrics from cloudwatch
func (s *server) MetricsGetImageHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)

	// get vars from the API route
	account := vars["account"]
	metric := vars["metric"]
	id := vars["id"]

	period := int64(300)
	// if p, ok := vars["period"]; ok {
	// 	// TODO p is a string, need int64
	// 	period = p
	// }

	start := "-P1D"
	if s, ok := vars["start"]; ok {
		start = s
	}

	end := "PT0H"
	if e, ok := vars["end"]; ok {
		end = e
	}

	cwService, ok := s.cloudwatchServices[account]
	if !ok {
		msg := fmt.Sprintf("cloudwatch service not found for account: %s", account)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}
	log.Debugf("found cloudwatch service %+v", cwService)

	metrics := []cloudwatch.Metric{
		cloudwatch.Metric{"AWS/EC2", metric, "InstanceId", id},
	}

	out, err := cwService.GetMetricWidget(r.Context(), metrics, period, start, end)
	if err != nil {
		log.Errorf("failed getting metrics widget image: %s", err)
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// MetricsGetURLHandler gets metrics from cloudwatch and returns a link to the image
func (s *server) MetricsGetImageUrlHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)

	// get vars from the API route
	account := vars["account"]
	metric := vars["metric"]
	id := vars["id"]

	period := int64(300)
	// 	// TODO p is a string, need int64
	// if p, ok := vars["period"]; ok {
	// 	period = p
	// }

	start := "-P1D"
	if s, ok := vars["start"]; ok {
		start = s
	}

	end := "PT0H"
	if e, ok := vars["end"]; ok {
		end = e
	}

	cwService, ok := s.cloudwatchServices[account]
	if !ok {
		msg := fmt.Sprintf("cloudwatch service not found for account: %s", account)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}
	log.Debugf("found cloudwatch service %+v", cwService)

	resultCache, ok := s.resultCache[account]
	if !ok {
		msg := fmt.Sprintf("result cache not found for account: %s", account)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}
	log.Debugf("found cost explorer result cache %+v", *resultCache)

	key := fmt.Sprintf("%s/%s/%s/%s/%s/%d", Org, id, metric, start, end, period)
	hashedCacheKey := s.imageCache.HashedKey(key)
	if res, expire, ok := resultCache.GetWithExpiration(hashedCacheKey); ok {
		log.Debugf("found cached object: %s", res)

		if body, ok := res.([]byte); ok {
			w.Header().Set("X-Cache-Hit", "true")
			w.Header().Set("X-Cache-Expire", fmt.Sprintf("%0.fs", time.Until(expire).Seconds()))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(body)
			return
		}
	}

	metrics := []cloudwatch.Metric{
		cloudwatch.Metric{"AWS/EC2", metric, "InstanceId", id},
	}

	image, err := cwService.GetMetricWidget(r.Context(), metrics, period, start, end)
	if err != nil {
		log.Errorf("failed getting metrics widget image: %s", err)
		handleError(w, err)
		return
	}

	meta, err := s.imageCache.Save(r.Context(), hashedCacheKey, image)
	if err != nil {
		log.Errorf("failed saving metrics widget image to cache: %s", err)
		handleError(w, err)
		return
	}
	resultCache.Set(hashedCacheKey, meta, 300*time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(meta)
}
