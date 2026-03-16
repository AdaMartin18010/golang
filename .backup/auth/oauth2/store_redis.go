package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisTokenStore Redis 令牌存储
type RedisTokenStore struct {
	client *redis.Client
	prefix string
}

// NewRedisTokenStore 创建 Redis 令牌存储
func NewRedisTokenStore(client *redis.Client, prefix string) (*RedisTokenStore, error) {
	if client == nil {
		return nil, errors.New("redis client is required")
	}

	if prefix == "" {
		prefix = "oauth2:"
	}

	return &RedisTokenStore{
		client: client,
		prefix: prefix,
	}, nil
}

// key 生成 Redis key
func (s *RedisTokenStore) key(k string) string {
	return s.prefix + "token:" + k
}

// refreshKey 生成刷新令牌 key
func (s *RedisTokenStore) refreshKey(refreshToken string) string {
	return s.prefix + "refresh:" + refreshToken
}

// Save 保存令牌
func (s *RedisTokenStore) Save(ctx context.Context, token *Token) error {
	// 序列化令牌
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// 保存访问令牌
	key := s.key(token.AccessToken)
	ttl := time.Until(token.ExpiresAt)
	if ttl > 0 {
		if err := s.client.Set(ctx, key, data, ttl).Err(); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
	} else {
		if err := s.client.Set(ctx, key, data, 0).Err(); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
	}

	// 保存刷新令牌索引
	if token.RefreshToken != "" {
		refreshKey := s.refreshKey(token.RefreshToken)
		if err := s.client.Set(ctx, refreshKey, token.AccessToken, ttl).Err(); err != nil {
			return fmt.Errorf("failed to save refresh token index: %w", err)
		}
	}

	return nil
}

// Get 获取令牌
func (s *RedisTokenStore) Get(ctx context.Context, accessToken string) (*Token, error) {
	key := s.key(accessToken)
	data, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// Delete 删除令牌
func (s *RedisTokenStore) Delete(ctx context.Context, accessToken string) error {
	// 获取令牌以找到刷新令牌
	token, err := s.Get(ctx, accessToken)
	if err == nil && token.RefreshToken != "" {
		// 删除刷新令牌索引
		refreshKey := s.refreshKey(token.RefreshToken)
		s.client.Del(ctx, refreshKey)
	}

	// 删除访问令牌
	key := s.key(accessToken)
	return s.client.Del(ctx, key).Err()
}

// DeleteByRefreshToken 通过刷新令牌删除
func (s *RedisTokenStore) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	// 获取访问令牌
	refreshKey := s.refreshKey(refreshToken)
	accessToken, err := s.client.Get(ctx, refreshKey).Result()
	if err == redis.Nil {
		return nil // 已不存在
	}
	if err != nil {
		return fmt.Errorf("failed to get access token from refresh token: %w", err)
	}

	// 删除访问令牌
	if err := s.Delete(ctx, accessToken); err != nil {
		return err
	}

	// 删除刷新令牌索引
	return s.client.Del(ctx, refreshKey).Err()
}

// CleanupExpired 清理过期令牌（Redis 自动过期，此方法用于手动清理）
func (s *RedisTokenStore) CleanupExpired(ctx context.Context) error {
	// Redis 使用 TTL 自动过期，不需要手动清理
	// 此方法保留用于兼容性
	return nil
}

// RedisClientStore Redis 客户端存储
type RedisClientStore struct {
	client *redis.Client
	prefix string
}

// NewRedisClientStore 创建 Redis 客户端存储
func NewRedisClientStore(client *redis.Client, prefix string) (*RedisClientStore, error) {
	if client == nil {
		return nil, errors.New("redis client is required")
	}

	if prefix == "" {
		prefix = "oauth2:"
	}

	return &RedisClientStore{
		client: client,
		prefix: prefix,
	}, nil
}

// key 生成 Redis key
func (s *RedisClientStore) key(clientID string) string {
	return s.prefix + "client:" + clientID
}

// Get 获取客户端
func (s *RedisClientStore) Get(ctx context.Context, clientID string) (*Client, error) {
	key := s.key(clientID)
	data, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrClientNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	var client Client
	if err := json.Unmarshal(data, &client); err != nil {
		return nil, fmt.Errorf("failed to unmarshal client: %w", err)
	}

	return &client, nil
}

// ValidateSecret 验证客户端密钥
func (s *RedisClientStore) ValidateSecret(ctx context.Context, clientID, secret string) error {
	client, err := s.Get(ctx, clientID)
	if err != nil {
		return err
	}

	if client.Secret != secret {
		return ErrInvalidClientSecret
	}

	return nil
}

// Save 保存客户端
func (s *RedisClientStore) Save(ctx context.Context, client *Client) error {
	data, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("failed to marshal client: %w", err)
	}

	key := s.key(client.ID)
	return s.client.Set(ctx, key, data, 0).Err() // 客户端不过期
}

// RedisCodeStore Redis 授权码存储
type RedisCodeStore struct {
	client *redis.Client
	prefix string
}

// NewRedisCodeStore 创建 Redis 授权码存储
func NewRedisCodeStore(client *redis.Client, prefix string) (*RedisCodeStore, error) {
	if client == nil {
		return nil, errors.New("redis client is required")
	}

	if prefix == "" {
		prefix = "oauth2:"
	}

	return &RedisCodeStore{
		client: client,
		prefix: prefix,
	}, nil
}

// key 生成 Redis key
func (s *RedisCodeStore) key(code string) string {
	return s.prefix + "code:" + code
}

// Save 保存授权码
func (s *RedisCodeStore) Save(ctx context.Context, code *AuthCode) error {
	data, err := json.Marshal(code)
	if err != nil {
		return fmt.Errorf("failed to marshal code: %w", err)
	}

	key := s.key(code.Code)
	ttl := time.Until(code.ExpiresAt)
	if ttl > 0 {
		return s.client.Set(ctx, key, data, ttl).Err()
	}
	return s.client.Set(ctx, key, data, 0).Err()
}

// Get 获取授权码
func (s *RedisCodeStore) Get(ctx context.Context, code string) (*AuthCode, error) {
	key := s.key(code)
	data, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrCodeNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get code: %w", err)
	}

	var authCode AuthCode
	if err := json.Unmarshal(data, &authCode); err != nil {
		return nil, fmt.Errorf("failed to unmarshal code: %w", err)
	}

	return &authCode, nil
}

// Delete 删除授权码
func (s *RedisCodeStore) Delete(ctx context.Context, code string) error {
	key := s.key(code)
	return s.client.Del(ctx, key).Err()
}

// CleanupExpired 清理过期授权码（Redis 自动过期）
func (s *RedisCodeStore) CleanupExpired(ctx context.Context) error {
	// Redis 使用 TTL 自动过期，不需要手动清理
	return nil
}
