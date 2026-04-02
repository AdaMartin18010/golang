package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"user-service/internal/domain"
)

type Cache interface {
	GetUser(key string) (*domain.User, error)
	SetUser(key string, user *domain.User, ttl time.Duration) error
	DeleteUser(key string) error
	Close() error
}

type redisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
		PoolSize: 10,
	})

	return &redisCache{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *redisCache) GetUser(key string) (*domain.User, error) {
	data, err := r.client.Get(r.ctx, "user:"+key).Result()
	if err == redis.Nil {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *redisCache) SetUser(key string, user *domain.User, ttl time.Duration) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, "user:"+key, data, ttl).Err()
}

func (r *redisCache) DeleteUser(key string) error {
	return r.client.Del(r.ctx, "user:"+key).Err()
}

func (r *redisCache) Close() error {
	return r.client.Close()
}
