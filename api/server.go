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

	log "github.com/sirupsen/logrus"
)

type server struct {
	router               *mux.Router
	version              common.Version
	context              context.Context
	costExplorerServices map[string]costexplorer.CostExplorer
}

// Org will carry throughout the api and get tagged on resources
var Org string

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
	}

	if config.Org == "" {
		return errors.New("'org' cannot be empty in the configuration")
	}
	Org = config.Org

	// Create a shared S3 session
	for name, c := range config.Accounts {
		log.Debugf("Creating new cost explorer service for account '%s' with key '%s' in region '%s' (org: %s)", name, c.Akid, c.Region, Org)
		s.costExplorerServices[name] = costexplorer.NewSession(c)
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

// rollBack executes functions from a stack of rollback functions
func rollBack(t *[]func() error) {
	if t == nil {
		return
	}

	tasks := *t
	log.Errorf("executing rollback of %d tasks", len(tasks))
	for i := len(tasks) - 1; i >= 0; i-- {
		f := tasks[i]
		if funcerr := f(); funcerr != nil {
			log.Errorf("rollback task error: %s, continuing rollback", funcerr)
		}
	}
}
