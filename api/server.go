package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/YaleSpinup/aws-go/services/session"
	"github.com/YaleSpinup/cost-api/common"
	"github.com/YaleSpinup/cost-api/imagecache"
	"github.com/YaleSpinup/cost-api/s3cache"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	cache "github.com/patrickmn/go-cache"

	log "github.com/sirupsen/logrus"
)

var (
	CacheExpireTime = 4 * time.Hour
	CachePurgeTime  = 15 * time.Minute
)

type server struct {
	accountsMap    map[string]string
	router         *mux.Router
	version        common.Version
	context        context.Context
	session        session.Session
	orgPolicy      string
	optimizerCache *cache.Cache
	resultCache    *cache.Cache
	imageCache     imagecache.ImageCache
	sessionCache   *cache.Cache
	org            string
}

// NewServer creates a new server and starts it
func NewServer(config common.Config) error {
	// setup server context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := server{
		accountsMap:  config.AccountsMap,
		router:       mux.NewRouter(),
		version:      config.Version,
		context:      ctx,
		sessionCache: cache.New(600*time.Second, 900*time.Second),
	}

	if config.Org == "" {
		return errors.New("'org' cannot be empty in the configuration")
	}
	s.org = config.Org

	orgPolicy, err := orgTagAccessPolicy(config.Org)
	if err != nil {
		return err
	}
	s.orgPolicy = orgPolicy

	if config.CacheExpireTime == "" {
		// set default expireTime
		log.Info("setting default cache expire time to 4h")
		config.CacheExpireTime = "4h"
	}

	exp, err := time.ParseDuration(config.CacheExpireTime)
	if err != nil {
		log.Error("Unexpected error with configured expiretime")
		return err
	}
	CacheExpireTime = exp

	if config.CachePurgeTime == "" {
		// set default purgeTime
		log.Info("setting default cache purge time to 15m")
		config.CachePurgeTime = "15m"
	}

	pt, err := time.ParseDuration(config.CachePurgeTime)
	if err != nil {
		log.Error("Unexpected error with configured purgetime")
		return err
	}
	CachePurgeTime = pt

	log.Debugf("creating new cost explorer result cache with expire time: %s and purge time: %s", CacheExpireTime, CachePurgeTime)
	s.resultCache = cache.New(CacheExpireTime, CachePurgeTime)

	log.Debugf("creating new optimizer cache with expire time: %s and purge time: %s", CacheExpireTime, CachePurgeTime)
	s.optimizerCache = cache.New(CacheExpireTime, CachePurgeTime)

	// Create a new session used for authentication and assuming cross account roles
	log.Debugf("Creating new session with key '%s' in region '%s'", config.Account.Akid, config.Account.Region)
	s.session = session.New(
		session.WithCredentials(config.Account.Akid, config.Account.Secret, ""),
		session.WithRegion(config.Account.Region),
		session.WithExternalID(config.Account.ExternalID),
		session.WithExternalRoleName(config.Account.Role),
	)

	// if specified, configure s3 image cache
	if config.ImageCache != nil {
		s.imageCache = s3cache.New(config.ImageCache)
	}

	publicURLs := map[string]string{
		"/v1/cost/ping":       "public",
		"/v1/cost/version":    "public",
		"/v1/cost/metrics":    "public",
		"/v1/metrics/ping":    "public",
		"/v1/metrics/version": "public",
		"/v1/metrics/metrics": "public",
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

// if we have an entry for the account name, return the associated account number
func (s *server) mapAccountNumber(name string) string {
	if a, ok := s.accountsMap[name]; ok {
		return a
	}
	return name
}
