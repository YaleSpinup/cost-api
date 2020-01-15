package api

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
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
	account := vars["account"]
	id := vars["id"]

	queries := r.URL.Query()
	metrics := queries["metric"]
	if len(metrics) == 0 {
		handleError(w, apierror.New(apierror.ErrBadRequest, "at least one metric is required", nil))
		return
	}

	period := int64(300)
	if p, ok := vars["period"]; ok && p != "" {
		dur, err := time.ParseDuration(p)
		if err != nil {
			msg := fmt.Sprintf("failed to parse period as duration: %s", err)
			handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
			return
		}

		period = int64(dur.Seconds())
	}

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

	cwMetrics := []cloudwatch.Metric{}
	for _, m := range metrics {
		cwMetrics = append(cwMetrics, cloudwatch.Metric{"AWS/EC2", m, "InstanceId", id})
	}

	out, err := cwService.GetMetricWidget(r.Context(), cwMetrics, period, start, end)
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
	account := vars["account"]
	id := vars["id"]

	queries := r.URL.Query()
	metrics := queries["metric"]
	if len(metrics) == 0 {
		handleError(w, apierror.New(apierror.ErrBadRequest, "at least one metric is required", nil))
		return
	}
	sort.Strings(metrics)

	period := int64(300)
	if p, ok := vars["period"]; ok && p != "" {
		dur, err := time.ParseDuration(p)
		if err != nil {
			msg := fmt.Sprintf("failed to parse period as duration: %s", err)
			handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
			return
		}

		period = int64(dur.Seconds())
	}

	start := "-P1D"
	if s, ok := vars["start"]; ok && s != "" {
		start = s
	}

	end := "PT0H"
	if e, ok := vars["end"]; ok && e != "" {
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

	key := fmt.Sprintf("%s/%s/%s/%s/%s/%d", Org, id, strings.Join(metrics, "-"), start, end, period)
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

	cwMetrics := []cloudwatch.Metric{}
	for _, m := range metrics {
		cwMetrics = append(cwMetrics, cloudwatch.Metric{"AWS/EC2", m, "InstanceId", id})
	}

	image, err := cwService.GetMetricWidget(r.Context(), cwMetrics, period, start, end)
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
