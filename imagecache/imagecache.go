package imagecache

import "context"

type ImageCache interface {
	GetMetadata(ctx context.Context, key string) ([]byte, error)
	Save(ctx context.Context, key string, obj []byte) ([]byte, error)
	HashedKey(key string) string
}
