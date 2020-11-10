package lock

import (
	"context"
	"sync"
	"time"

	"github.com/bsm/redislock"
	redisCli "github.com/go-redis/redis/v8"
)

var (
	redisOnce = sync.Once{}
	redis     *Redis
)

type Redis struct {
	cli *redisCli.Client
}

type RedisConfig interface {
	Host() string
	Pwd() string
}

type Lock interface {
	Unlock() error
}

type RemoteStorage interface {
	ObtainLock(key string, duration time.Duration) (Lock, error)
}

type RedisLock struct {
	lock *redislock.Lock
}

func (l *RedisLock) Unlock() error {
	return l.lock.Release(context.Background())
}

func NewRedisStorage(cfg RedisConfig) RemoteStorage {
	redisOnce.Do(func() {
		redis = &Redis{
			cli: redisCli.NewClient(&redisCli.Options{
				Addr:     cfg.Host(),
				Password: cfg.Pwd(),
			}),
		}
	})
	return redis
}

func (r *Redis) ObtainLock(key string, duration time.Duration) (Lock, error) {
	var (
		locker = redislock.New(r.cli)
	)
	lock, err := locker.Obtain(context.Background(), key, duration, nil)

	if err != nil {
		return nil, err
	}

	return &RedisLock{
		lock: lock,
	}, nil
}
