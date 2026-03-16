// +build redis

package middleware

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisAdapter 将 go-redis/v9 客户端适配为 RedisClient 接口
type RedisAdapter struct {
	client *redis.Client
}

// NewRedisAdapter 创建 Redis 适配器
func NewRedisAdapter(client *redis.Client) RedisClient {
	return &RedisAdapter{client: client}
}

// Pipeline 创建 Pipeline
func (a *RedisAdapter) Pipeline() RedisPipeline {
	return &RedisPipelineAdapter{pipe: a.client.Pipeline()}
}

// ZRemRangeByScore 移除分数范围内的成员
func (a *RedisAdapter) ZRemRangeByScore(ctx context.Context, key, min, max string) *RedisIntCmd {
	cmd := a.client.ZRemRangeByScore(ctx, key, min, max)
	return &RedisIntCmd{val: cmd.Val(), err: cmd.Err()}
}

// ZCard 获取有序集合的成员数
func (a *RedisAdapter) ZCard(ctx context.Context, key string) *RedisIntCmd {
	cmd := a.client.ZCard(ctx, key)
	return &RedisIntCmd{val: cmd.Val(), err: cmd.Err()}
}

// ZAdd 添加有序集合成员
func (a *RedisAdapter) ZAdd(ctx context.Context, key string, members ...interface{}) *RedisIntCmd {
	cmd := a.client.ZAdd(ctx, key, members...)
	return &RedisIntCmd{val: cmd.Val(), err: cmd.Err()}
}

// Expire 设置过期时间
func (a *RedisAdapter) Expire(ctx context.Context, key string, expiration time.Duration) *RedisBoolCmd {
	cmd := a.client.Expire(ctx, key, expiration)
	return &RedisBoolCmd{val: cmd.Val(), err: cmd.Err()}
}

// Exec 执行 Pipeline
func (a *RedisAdapter) Exec(ctx context.Context) ([]RedisCmder, error) {
	cmds, err := a.client.Pipeline().Exec(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]RedisCmder, len(cmds))
	for i, cmd := range cmds {
		result[i] = &RedisCmdAdapter{cmd: cmd}
	}
	return result, nil
}

// RedisPipelineAdapter Pipeline 适配器
type RedisPipelineAdapter struct {
	pipe redis.Pipeliner
}

// ZRemRangeByScore 移除分数范围内的成员
func (p *RedisPipelineAdapter) ZRemRangeByScore(ctx context.Context, key, min, max string) *RedisIntCmd {
	cmd := p.pipe.ZRemRangeByScore(ctx, key, min, max)
	return &RedisIntCmd{err: cmd.Err()}
}

// ZCard 获取有序集合的成员数
func (p *RedisPipelineAdapter) ZCard(ctx context.Context, key string) *RedisIntCmd {
	cmd := p.pipe.ZCard(ctx, key)
	return &RedisIntCmd{err: cmd.Err()}
}

// ZAdd 添加有序集合成员
func (p *RedisPipelineAdapter) ZAdd(ctx context.Context, key string, members ...interface{}) *RedisIntCmd {
	cmd := p.pipe.ZAdd(ctx, key, members...)
	return &RedisIntCmd{err: cmd.Err()}
}

// Expire 设置过期时间
func (p *RedisPipelineAdapter) Expire(ctx context.Context, key string, expiration time.Duration) *RedisBoolCmd {
	cmd := p.pipe.Expire(ctx, key, expiration)
	return &RedisBoolCmd{err: cmd.Err()}
}

// Exec 执行 Pipeline
func (p *RedisPipelineAdapter) Exec(ctx context.Context) ([]RedisCmder, error) {
	cmds, err := p.pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]RedisCmder, len(cmds))
	for i, cmd := range cmds {
		result[i] = &RedisCmdAdapter{cmd: cmd}
	}
	return result, nil
}

// RedisCmdAdapter 命令适配器
type RedisCmdAdapter struct {
	cmd redis.Cmder
}

// Err 返回错误
func (c *RedisCmdAdapter) Err() error {
	return c.cmd.Err()
}
