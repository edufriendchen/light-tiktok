package cache

import "context"

type (
	Client interface {
		LockClient
		Ping(ctx context.Context) error
		Set(ctx context.Context, key string, value interface{}, opts ...Option) error
		Get(ctx context.Context, key string, result interface{}) error
	}

	LockClient interface {
		Lock(ctx context.Context, key string, opts ...Option) (LockClient, error)
		Unlock(ctx context.Context) error
	}
)
