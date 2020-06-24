package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/YaleSpinup/cost-api/cloudwatch"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// GetEC2MetricsURLHandler gets metrics from cloudwatch and returns a link to the image
func (s *server) GetEC2MetricsURLHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	instanceId := vars["id"]

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

	queries := r.URL.Query()
	metrics := queries["metric"]
	if len(metrics) == 0 {
		handleError(w, apierror.New(apierror.ErrBadRequest, "at least one metric is required", nil))
		return
	}

	req := cloudwatch.MetricsRequest{}
	if err := parseQuery(r, req); err != nil {
		handleError(w, apierror.New(apierror.ErrBadRequest, "failed to parse query", err))
		return
	}

	key := fmt.Sprintf("%s/%s/%s%s", Org, instanceId, strings.Join(metrics, "-"), req.String())
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
		cwMetrics = append(cwMetrics, cloudwatch.Metric{"AWS/EC2", m, "InstanceId", instanceId})
	}
	req["metrics"] = cwMetrics

	log.Debugf("getting metrics with request %+v", req)
	image, err := cwService.GetMetricWidget(r.Context(), req)
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

// GetECSMetricsURLHandler gets metrics from cloudwatch and returns a link to the image
func (s *server) GetECSMetricsURLHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	cluster := vars["cluster"]
	service := vars["service"]

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

	queries := r.URL.Query()
	metrics := queries["metric"]
	if len(metrics) == 0 {
		handleError(w, apierror.New(apierror.ErrBadRequest, "at least one metric is required", nil))
		return
	}

	req := cloudwatch.MetricsRequest{}
	if err := parseQuery(r, req); err != nil {
		handleError(w, apierror.New(apierror.ErrBadRequest, "failed to parse query", err))
		return
	}

	key := fmt.Sprintf("%s/%s/%s%s", Org, fmt.Sprintf("%s-%s", cluster, service), strings.Join(metrics, "-"), req.String())
	log.Debugf("object key: %s", key)

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
		cwMetrics = append(cwMetrics, cloudwatch.Metric{"AWS/ECS", m, "ClusterName", cluster, "ServiceName", service})
	}
	req["metrics"] = cwMetrics

	log.Debugf("getting metrics with request %+v", req)
	image, err := cwService.GetMetricWidget(r.Context(), req)
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

// GetS3MetricsURLHandler gets metrics from cloudwatch and returns a link to the image
func (s *server) GetS3MetricsURLHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	bucketName := vars["bucket"]
	metric := vars["metric"]

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

	// only support NumberOfObjects and BucketSizeBytes
	var storageType string
	switch metric {
	case "BucketSizeBytes":
		storageType = "StandardStorage"
	case "NumberOfObjects":
		storageType = "AllStorageTypes"
	default:
		msg := fmt.Sprintf("invalid metric requested: %s", metric)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, nil))
		return
	}

	req := cloudwatch.MetricsRequest{
		"period": int64(86400),
		"stat":   "Maximum",
		"start":  "-P30D",
		"end":    "PT0H",
		"metrics": []cloudwatch.Metric{
			{"AWS/S3", metric, "StorageType", storageType, "BucketName", bucketName},
		},
	}

	key := fmt.Sprintf("%s/%s/%s%s", Org, bucketName, metric, req.String())
	log.Debugf("object key: %s", key)

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

	log.Debugf("getting metrics with request %+v", req)
	image, err := cwService.GetMetricWidget(r.Context(), req)
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

// GetRDSMetricsURLHandler gets RDS metrics from cloudwatch and returns a link to the image
// Can use DBInstanceIdentifier or DBClusterIdentifier
func (s *server) GetRDSMetricsURLHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	queryType := vars["type"]
	instanceId := vars["id"]

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

	queries := r.URL.Query()
	metrics := queries["metric"]
	if len(metrics) == 0 {
		handleError(w, apierror.New(apierror.ErrBadRequest, "at least one metric is required", nil))
		return
	}

	req := cloudwatch.MetricsRequest{}
	if err := parseQuery(r, req); err != nil {
		handleError(w, apierror.New(apierror.ErrBadRequest, "failed to parse query", err))
		return
	}

	key := fmt.Sprintf("%s/%s/%s%s", Org, instanceId, strings.Join(metrics, "-"), req.String())
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
		switch queryType {
		case "instance":
			cwMetrics = append(cwMetrics, cloudwatch.Metric{"AWS/RDS", m, "DBInstanceIdentifier", instanceId})
		case "cluster":
			cwMetrics = append(cwMetrics, cloudwatch.Metric{"AWS/RDS", m, "DBClusterIdentifier", instanceId})
		default:
			msg := fmt.Sprintf("invalid type requested: %s", m)
			handleError(w, apierror.New(apierror.ErrBadRequest, msg, nil))
			return
		}
	}
	req["metrics"] = cwMetrics

	log.Debugf("getting metrics with request %+v", req)
	image, err := cwService.GetMetricWidget(r.Context(), req)
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

func parseQuery(r *http.Request, request cloudwatch.MetricsRequest) error {
	log.SetLevel(log.DebugLevel)
	queries := r.URL.Query()
	log.Debugf("parsing queries: %+v", queries)

	stat := "Average"
	if s, ok := queries["stat"]; ok {
		stat = s[0]
	}
	request["stat"] = stat

	period := int64(300)
	if p, ok := queries["period"]; ok && p[0] != "" {
		dur, err := time.ParseDuration(p[0])
		if err != nil {
			return errors.Wrap(err, "failed to parse period as duration")
		}

		period = int64(dur.Seconds())
	}
	request["period"] = period

	start := "-P1D"
	if s, ok := queries["start"]; ok {
		start = s[0]
	}
	request["start"] = start

	end := "PT0H"
	if e, ok := queries["end"]; ok {
		end = e[0]
	}
	request["end"] = end

	height := int64(400)
	if h, ok := queries["height"]; ok {
		hint, err := strconv.ParseInt(h[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse height as int")
		}

		if hint < int64(1) || hint > int64(2000) {
			return fmt.Errorf("invalid height %d, value must be >=1 and <= 2000", hint)
		}

		height = hint
	}
	request["height"] = height

	width := int64(600)
	if w, ok := queries["width"]; ok {
		wint, err := strconv.ParseInt(w[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse width as int")
		}

		if wint < int64(1) || wint > int64(2000) {
			return fmt.Errorf("invalid width %d, value must be >=1 and <= 2000", wint)
		}

		width = wint
	}
	request["width"] = width

	return nil
}
