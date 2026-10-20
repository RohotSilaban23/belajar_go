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

func (r *RedisService) IsTokenActive(
	ctx context.Context,
	token string,
) (bool, error) {

	key := "active_token:" + token

	exists, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (r *RedisService) SaveOTP(
	ctx context.Context,
	email string,
	otpHash string,
	expiration time.Duration,
) error {
	key := "forgot_password_otp:" + email
	return r.Client.Set(ctx, key, otpHash, expiration).Err()
}

func (r *RedisService) GetOTP(
	ctx context.Context,
	email string,
) (string, error) {
	key := "forgot_password_otp:" + email
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisService) DeleteOTP(
	ctx context.Context,
	email string,
) error {
	key := "forgot_password_otp:" + email
	return r.Client.Del(ctx, key).Err()
}
