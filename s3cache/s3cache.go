package s3cache

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fossoreslp/go-uuid-v4"

	"github.com/YaleSpinup/cost-api/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	log "github.com/sirupsen/logrus"
)

type S3Cache struct {
	Service             s3iface.S3API
	Bucket              string
	Prefix              string
	LoggingBucket       string
	LoggingBucketPrefix string
	HashingToken        string
}

// New creates a new S3 session and adds some config data
func New(s3cache *common.S3Cache) *S3Cache {
	log.Infof("creating new aws session for S3 with key id %s in region %s", s3cache.Akid, s3cache.Region)

	s := S3Cache{}
	config := aws.Config{
		Credentials: credentials.NewStaticCredentials(s3cache.Akid, s3cache.Secret, ""),
		Region:      aws.String(s3cache.Region),
	}

	if s3cache.Endpoint != "" {
		config.Endpoint = aws.String(s3cache.Endpoint)
	}

	if s3cache.Bucket == "" {
		log.Error("s3 cache bucket name is required")
		return nil
	}
	s.Bucket = s3cache.Bucket
	s.Prefix = s3cache.Prefix

	sess := session.Must(session.NewSession(&config))
	s.Service = s3.New(sess)
	if s3cache.AccessLog != nil {
		s.LoggingBucket = s3cache.AccessLog.Bucket
		s.LoggingBucketPrefix = s3cache.AccessLog.Prefix
	}

	// if a static hasing token is passed, use it, otherwise generate one
	if s3cache.HashingToken != "" {
		s.HashingToken = s3cache.HashingToken
	} else {
		uuidv4, _ := uuid.New()
		s.HashingToken = uuidv4.String()
	}

	return &s
}

func (s *S3Cache) GetMetadata(ctx context.Context, key string) ([]byte, error) {
	log.Infof("getting %s metadata from cache", key)

	if s.Prefix != "" {
		key = s.Prefix + "/" + key
	}

	obj, err := s.Service.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		msg := fmt.Sprintf("failed to find object %s in bucket %s: %s", key, s.Bucket, err)
		return nil, ErrCode(msg, err)
	}

	log.Infof("got response from s3: %+v", obj)

	// TODO: this only works in us-east-1, other regions are "https://s3-<region-id>.amazonaws.com/<bucket>/<key>"
	resp := struct {
		ImageURL string
	}{
		ImageURL: fmt.Sprintf("https://s3.amazonaws.com/%s/%s", s.Bucket, key),
	}

	return json.Marshal(resp)
}

func (s *S3Cache) Save(ctx context.Context, key string, obj []byte) ([]byte, error) {
	log.Infof("saving object %s to cache", key)

	if s.Prefix != "" {
		key = s.Prefix + "/" + key
	}

	_, err := s.Service.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(obj),
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		msg := fmt.Sprintf("error saving object %s to bucket %s: %s", key, s.Bucket, err)
		return nil, ErrCode(msg, err)
	}

	// TODO: this only works in us-east-1, other regions are "https://s3-<region-id>.amazonaws.com/<bucket>/<key>"
	resp := struct {
		ImageURL string
	}{
		ImageURL: fmt.Sprintf("https://s3.amazonaws.com/%s/%s", s.Bucket, key),
	}

	return json.Marshal(resp)
}

func (s *S3Cache) HashedKey(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s.HashingToken + key))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
