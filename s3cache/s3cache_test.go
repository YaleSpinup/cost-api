package s3cache

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/cost-api/common"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// mockS3Client is a fake costexplorer client
type mockS3Client struct {
	s3iface.S3API
	t   *testing.T
	err error
}

func newmockS3Client(t *testing.T, err error) s3iface.S3API {
	return &mockS3Client{
		t:   t,
		err: err,
	}
}

func (m *mockS3Client) HeadObjectWithContext(context.Context, *s3.HeadObjectInput, ...request.Option) (*s3.HeadObjectOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	return nil, nil
}

func (m *mockS3Client) PutObjectWithContext(context.Context, *s3.PutObjectInput, ...request.Option) (*s3.PutObjectOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	return nil, nil
}

func TestNewSession(t *testing.T) {
	e := New(&common.S3Cache{})
	if to := reflect.TypeOf(e).String(); to != "*s3cache.S3Cache" {
		t.Errorf("expected type to be '*s3cache.S3Cache', got %s", to)
	}
}

func TestGetMetadata(t *testing.T) {
	c := &S3Cache{
		Service:      newmockS3Client(t, nil),
		Bucket:       "testbucket",
		Prefix:       "foobar",
		HashingToken: "test",
	}

	expected := "{\"ImageURL\":\"https://s3.amazonaws.com/testbucket/foobar/somekey\"}"
	out, err := c.GetMetadata(context.TODO(), "somekey")
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	if string(out) != expected {
		t.Errorf("expected %s, got %s", expected, string(out))
	}

	c.Service.(*mockS3Client).err = awserr.New(s3.ErrCodeNoSuchKey, "no such key", nil)
	_, err = c.GetMetadata(context.TODO(), "somekey")
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrNotFound {
			t.Errorf("expected error code %s, got: %s", apierror.ErrNotFound, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}

	// test non-aws error
	c.Service.(*mockS3Client).err = errors.New("things blowing up!")
	_, err = c.GetMetadata(context.TODO(), "somekey")
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrInternalError {
			t.Errorf("expected error code %s, got: %s", apierror.ErrInternalError, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}
}

func TestSave(t *testing.T) {
	c := &S3Cache{
		Service:      newmockS3Client(t, nil),
		Bucket:       "testbucket",
		Prefix:       "foobar",
		HashingToken: "test",
	}

	expected := "{\"ImageURL\":\"https://s3.amazonaws.com/testbucket/foobar/somekey\"}"
	out, err := c.Save(context.TODO(), "somekey", []byte{})
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	if string(out) != expected {
		t.Errorf("expected %s, got %s", expected, string(out))
	}

	c.Service.(*mockS3Client).err = awserr.New(s3.ErrCodeNoSuchBucket, "no such bucket", nil)
	_, err = c.Save(context.TODO(), "somekey", []byte{})
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrNotFound {
			t.Errorf("expected error code %s, got: %s", apierror.ErrNotFound, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}

	// test non-aws error
	c.Service.(*mockS3Client).err = errors.New("things blowing up!")
	_, err = c.Save(context.TODO(), "somekey", []byte{})
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrInternalError {
			t.Errorf("expected error code %s, got: %s", apierror.ErrInternalError, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}
}

func TestHashedKey(t *testing.T) {
	c := &S3Cache{
		Service:      newmockS3Client(t, nil),
		Bucket:       "testbucket",
		Prefix:       "foobar",
		HashingToken: "test",
	}

	expected := "FL4euODSAlXXPEEekKy-YwsT6bnkP6rVJJg0FTEc07c="
	out := c.HashedKey("somekey")
	if out != expected {
		t.Errorf("expected %s, got %s", expected, out)
	}
}
