package cacheadapter

import (
	"context"
	"time"

	domaincache "github.com/mohamedfawas/employee_management_system/internal/domain/cache"
	pkgredis "github.com/mohamedfawas/employee_management_system/pkg/cache"
)

type RedisAdapter struct {
	client *pkgredis.Client
}

func NewRedisAdapter(client *pkgredis.Client) domaincache.Cache {
	return &RedisAdapter{client: client}
}

func (r *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key)
}

func (r *RedisAdapter) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl)
}

func (r *RedisAdapter) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key)
}

func (r *RedisAdapter) Exists(ctx context.Context, key string) (bool, error) {
	return r.client.Exists(ctx, key)
}
