package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YaleSpinup/apierror"
	ce "github.com/YaleSpinup/cost-api/costexplorer"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// PingHandler responds to ping requests
func (s *server) PingHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	log.Debug("Ping/Pong")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

// VersionHandler responds to version requests
func (s *server) VersionHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(struct {
		Version    string `json:"version"`
		GitHash    string `json:"githash"`
		BuildStamp string `json:"buildstamp"`
	}{
		Version:    fmt.Sprintf("%s%s", s.version.Version, s.version.VersionPrerelease),
		GitHash:    s.version.GitHash,
		BuildStamp: s.version.BuildStamp,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleError handles standard apierror return codes
func handleError(w http.ResponseWriter, err error) {
	log.Error(err.Error())
	if aerr, ok := errors.Cause(err).(apierror.Error); ok {
		switch aerr.Code {
		case apierror.ErrForbidden:
			w.WriteHeader(http.StatusForbidden)
		case apierror.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		case apierror.ErrConflict:
			w.WriteHeader(http.StatusConflict)
		case apierror.ErrBadRequest:
			w.WriteHeader(http.StatusBadRequest)
		case apierror.ErrLimitExceeded:
			w.WriteHeader(http.StatusTooManyRequests)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(aerr.String()))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
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
