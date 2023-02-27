package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	redis_store "github.com/eko/gocache/store/redis/v4"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.opencensus.io/resource"
)

type client struct {
	resource     resource.Resource
	db           *redis.Client
	redSync      *redsync.Redsync
	mutex        *redsync.Mutex
	cacheManager *cache.ChainCache[any]
}

func (c *client) Lock(ctx context.Context, key string, opts ...Option) (LockClient, error) {

	cacheOptions := &options{}
	for _, opt := range opts {
		opt.Apply(cacheOptions)
	}

	var lockOptions []redsync.Option
	if cacheOptions.expiry > 0 {
		lockOptions = append(lockOptions, redsync.WithExpiry(cacheOptions.expiry))
	}

	mutex := c.redSync.NewMutex(key, lockOptions...)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, err
	}

	return &client{
		resource: c.resource,
		db:       c.db,
		redSync:  c.redSync,
		mutex:    mutex,
	}, nil
}

func (c client) Unlock(ctx context.Context) error {
	if c.mutex == nil {
		return errors.New("ErrNothingToUnlock")
	}
	if ok, err := c.mutex.UnlockContext(ctx); !ok || err != nil {
		if err == nil {
			err = errors.New("ErrFailedToUnlockRedis")
		}
		return err
	}

	return nil
}

func (c client) Ping(ctx context.Context) error {
	return c.db.Ping(ctx).Err()
}

func (c client) Set(ctx context.Context, key string, value interface{}, opts ...Option) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.cacheManager.Set(ctx, key, json)
	if err != nil {
		return err
	}
	fmt.Println("SET cache key:", key, "value:", value)
	return nil
}

func (c client) Get(ctx context.Context, key string, result interface{}) error {
	value, err := c.cacheManager.Get(ctx, key)
	if err == nil {
		str := fmt.Sprintf("%s", value)
		fmt.Println("value-str:", str)
		if err = json.Unmarshal([]byte(str), &result); err == nil {
			fmt.Println("GET cache key:", key, "value:", &result)
			return nil
		}
	}
	return err
}

func NewClient(resource resource.Resource) (Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "1.117.77.53:6379",
		Password: "friendchen",
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Println("error when connect to redis: " + err.Error())
		return nil, err
	}

	pool := goredis.NewPool(redisClient)
	rs := redsync.New(pool)

	redisStore := redis_store.NewRedis(redis.NewClient(&redis.Options{
		Addr:     "1.117.77.53:6379",
		Password: "friendchen",
	}))
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{NumCounters: 1000, MaxCost: 100, BufferItems: 64})
	if err != nil {
		panic(err)
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)

	// Initialize chained cache
	cacheManager := cache.NewChain[any](
		cache.New[any](ristrettoStore),
		cache.New[any](redisStore),
	)

	return &client{
		resource:     resource,
		db:           redisClient,
		redSync:      rs,
		cacheManager: cacheManager,
	}, nil
}
