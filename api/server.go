package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/YaleSpinup/cost-api/cloudwatch"
	"github.com/YaleSpinup/cost-api/common"
	"github.com/YaleSpinup/cost-api/costexplorer"
	"github.com/YaleSpinup/cost-api/imagecache"
	"github.com/YaleSpinup/cost-api/s3cache"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	cache "github.com/patrickmn/go-cache"

	log "github.com/sirupsen/logrus"
)

var (
	// Org will carry throughout the api and get tagged on resources
	Org string
)

type server struct {
	router               *mux.Router
	version              common.Version
	context              context.Context
	costExplorerServices map[string]costexplorer.CostExplorer
	cloudwatchServices   map[string]cloudwatch.Cloudwatch
	resultCache          map[string]*cache.Cache
	imageCache           imagecache.ImageCache
}

// NewServer creates a new server and starts it
func NewServer(config common.Config) error {
	// setup server context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := server{
		router:               mux.NewRouter(),
		version:              config.Version,
		context:              ctx,
		costExplorerServices: make(map[string]costexplorer.CostExplorer),
		cloudwatchServices:   make(map[string]cloudwatch.Cloudwatch),
		resultCache:          make(map[string]*cache.Cache),
	}

	if config.Org == "" {
		return errors.New("'org' cannot be empty in the configuration")
	}
	Org = config.Org

	if config.CacheExpireTime == "" {
		// set default expireTime
		log.Info("setting default cache expire time to 4h")
		config.CacheExpireTime = "4h"
	}

	expireTime, err := time.ParseDuration(config.CacheExpireTime)
	if err != nil {
		log.Error("Unexpected error with configured expiretime")
		return err
	}

	if config.CachePurgeTime == "" {
		// set default purgeTime
		log.Info("setting default cache purge time to 15m")
		config.CachePurgeTime = "15m"
	}

	purgeTime, err := time.ParseDuration(config.CachePurgeTime)
	if err != nil {
		log.Error("Unexpected error with configured purgetime")
		return err
	}

	// Create shared cost explorer sessions, cloudwatch sessions, and go-cache instances per account defined in the config
	for name, c := range config.Accounts {
		log.Debugf("creating new cost explorer service for account '%s' with key '%s' in region '%s' (org: %s)", name, c.Akid, c.Region, Org)
		s.costExplorerServices[name] = costexplorer.NewSession(c)

		log.Debugf("creating new cloudwatch service for account '%s' with key '%s' in region '%s' (org: %s)", name, c.Akid, c.Region, Org)
		s.cloudwatchServices[name] = cloudwatch.NewSession(c)

		log.Debugf("creating new result cache for account '%s' with expire time: %s and purge time: %s", name, expireTime.String(), purgeTime.String())
		s.resultCache[name] = cache.New(expireTime, purgeTime)
	}

	// if specified, configure s3 image cache
	if config.ImageCache != nil {
		s.imageCache = s3cache.New(config.ImageCache)
	}

	publicURLs := map[string]string{
		"/v1/cost/ping":    "public",
		"/v1/cost/version": "public",
		"/v1/cost/metrics": "public",
	}

	// load routes
	s.routes()

	if config.ListenAddress == "" {
		config.ListenAddress = ":8080"
	}
	handler := handlers.RecoveryHandler()(handlers.LoggingHandler(os.Stdout, TokenMiddleware([]byte(config.Token), publicURLs, s.router)))
	srv := &http.Server{
		Handler:      handler,
		Addr:         config.ListenAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Starting listener on %s", config.ListenAddress)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// LogWriter is an http.ResponseWriter
type LogWriter struct {
	http.ResponseWriter
}

// Write log message if http response writer returns an error
func (w LogWriter) Write(p []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(p)
	if err != nil {
		log.Errorf("Write failed: %v", err)
	}
	return
}
