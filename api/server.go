package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/YaleSpinup/cost-api/common"
	"github.com/YaleSpinup/cost-api/costexplorer"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	cache "github.com/patrickmn/go-cache"

	log "github.com/sirupsen/logrus"
)

var (
	// Org will carry throughout the api and get tagged on resources
	Org string

	// ResultsCache is a map of in-memory caches
	// ResultsCache = make(map[string]*cache.Cache)
)

type server struct {
	router               *mux.Router
	version              common.Version
	context              context.Context
	costExplorerServices map[string]costexplorer.CostExplorer
	resultCache          map[string]*cache.Cache
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
		resultCache:          make(map[string]*cache.Cache),
	}

	if config.Org == "" {
		return errors.New("'org' cannot be empty in the configuration")
	}
	Org = config.Org

	if config.CacheExpireTime == "" {
		// set default expireTime
		config.CacheExpireTime = "4h"
	}

	expireTime, err := time.ParseDuration(config.CacheExpireTime)
	if err != nil {
		log.Warn("Unexpected error with configured expiretime")
	}

	if config.CachePurgeTime == "" {
		// set default purgeTime
		config.CachePurgeTime = "15m"
	}

	purgeTime, err := time.ParseDuration(config.CachePurgeTime)
	if err != nil {
		log.Warn("Unexpected error with configured purgetime")
	}

	log.Debugf("default cache expire time is: %s", cache.DefaultExpiration.String())

	// Create a shared Cost Explorer session
	for name, c := range config.Accounts {
		log.Debugf("Creating new cost explorer service for account '%s' with key '%s' in region '%s' (org: %s)", name, c.Akid, c.Region, Org)
		s.costExplorerServices[name] = costexplorer.NewSession(c)

		// Create a cache with a 4 hour default expiry and a 15 minute default purge time
		// ResultsCache[name] = cache.New(expireTime, purgeTime)
		s.resultCache[name] = cache.New(expireTime, purgeTime)
	}

	log.Debugf("default cache expire time is: %s", cache.DefaultExpiration.String())

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
	handler := handlers.RecoveryHandler()(handlers.LoggingHandler(os.Stdout, TokenMiddleware(config.Token, publicURLs, s.router)))
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
