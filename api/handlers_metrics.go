package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/YaleSpinup/cost-api/apierror"
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

	cwService, ok := s.cloudwatchServices[account]
	if !ok {
		msg := fmt.Sprintf("cloudwatch service not found for account: %s", account)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}
	log.Debugf("found cloudwatch service %+v", cwService)

	out, err := cwService.GetMetricWidget(r.Context(), metric, id)
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

	hashedCacheKey := s.imageCache.HashedKey(Org + "/" + id + "/" + metric)
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

	image, err := cwService.GetMetricWidget(r.Context(), metric, id)
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
