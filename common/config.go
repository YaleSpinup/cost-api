package common

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Config is representation of the configuration data
type Config struct {
	ListenAddress   string
	Accounts        map[string]Account
	Account         Account
	Token           string
	LogLevel        string
	Version         Version
	Org             string
	CacheExpireTime string
	CachePurgeTime  string
	ImageCache      *S3Cache
}

// Account is the configuration for an individual account
type Account struct {
	Akid       string
	Endpoint   string
	ExternalID string
	Region     string
	Role       string
	Secret     string
}

// Version carries around the API version information
type Version struct {
	Version           string
	VersionPrerelease string
	BuildStamp        string
	GitHash           string
}

type S3Cache struct {
	Bucket       string
	Endpoint     string
	Region       string
	Akid         string
	Secret       string
	Prefix       string
	HashingToken string
	AccessLog    *AccessLog
}

// AccessLog is the configuration for a bucket's access log
type AccessLog struct {
	Bucket string
	Prefix string
}

// ReadConfig decodes the configuration from an io Reader
func ReadConfig(r io.Reader) (Config, error) {
	var c Config
	log.Infoln("Reading configuration")
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return c, errors.Wrap(err, "unable to decode JSON message")
	}
	return c, nil
}
