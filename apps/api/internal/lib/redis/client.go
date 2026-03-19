package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(redisUrl string, password string) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: password,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{
		Client: rdb,
	}, nil
}

func (r *RedisClient) Close() {
	r.Client.Close()
}

func (r *RedisClient) Health(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
