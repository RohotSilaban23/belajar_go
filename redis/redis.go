package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client *goredis.Client
}

func NewRedisService() *RedisService {
	client := goredis.NewClient(&goredis.Options{
		Addr: "localhost:6379",
	})

	return &RedisService{
		Client: client,
	}
}

func (r *RedisService) SaveToken(
	ctx context.Context,
	token string,
	userID int64,
	expiration time.Duration,
) error {

	key := "active_token:" + token

	return r.Client.Set(
		ctx,
		key,
		userID,
		expiration,
	).Err()
}

func (r *RedisService) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *RedisService) DeleteToken(
	ctx context.Context,
	token string,
) error {

	key := "active_token:" + token

	return r.Client.Del(ctx, key).Err()
}
