package cache

import (
	"context"
	"time"

	"github.com/jinmukeji/huimaibao-service/pkg/cache/store"
)

// CacheInterface represents the interface for all caches (aggregates, metric, memory, redis, ...)
type CacheInterface interface {
	// Get 通过key获取内容，找不到不会报错，返回nil
	Get(ctx context.Context, key string) (interface{}, error)
	// GetMulti 获取多个key的内容，找不到不会报错，返回nil
	GetMulti(ctx context.Context, keys []string) ([]interface{}, error)
	// Put 设置key，value
	// 更新就是直接重新put即可
	Put(ctx context.Context, key string, val interface{}, opts ...store.Option) error
	// Incr value加一
	Incr(ctx context.Context, key string) error
	// Decr value减一
	Decr(ctx context.Context, key string) error
	// Delete 删除key
	Delete(ctx context.Context, key string) error
	// IsExist 检测key是否存在
	IsExist(ctx context.Context, key string) (bool, error)
	// ClearByTags 通过tag批量删除
	ClearByTags(ctx context.Context, tags []string, opts ...store.Option) error
	// HSet 设置map
	HSet(ctx context.Context, key string, val map[string]interface{}, opts ...store.Option) error
	// HGet 通过key和map的key获取map指定字段值
	HGet(ctx context.Context, key, field string) (interface{}, error)
	// HGetAll 返回所有map字段和对应的值
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	// HIncrBy map的field对应value加上一个加数
	HIncrBy(ctx context.Context, key, field string, incr int64) error
	// Expire 设置过期时间
	Expire(ctx context.Context, key string, duration time.Duration) error
}

// Cache represents the configuration needed by a cache
type Cache struct {
	store store.StoreInterface
}

// New instantiates a new cache entry
func New(store store.StoreInterface) *Cache {
	return &Cache{
		store: store,
	}
}

// Get 通过key获取内容，找不到不会报错，返回nil
func (c *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	return c.store.Get(ctx, key)
}

// GetMulti 获取多个key的内容，找不到不会报错，返回nil
func (c *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	return c.store.GetMulti(ctx, keys)
}

// Put 设置key，value
// 更新就是直接重新put即可
func (c *Cache) Put(ctx context.Context, key string, val interface{}, opts ...store.Option) error {
	return c.store.Put(ctx, key, val, opts...)
}

// Incr value加一
func (c *Cache) Incr(ctx context.Context, key string) error {
	return c.store.Incr(ctx, key)
}

// Decr value减一
func (c *Cache) Decr(ctx context.Context, key string) error {
	return c.store.Decr(ctx, key)
}

// Delete 删除key
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.store.Delete(ctx, key)
}

// IsExist 检测key是否存在
func (c *Cache) IsExist(ctx context.Context, key string) (bool, error) {
	return c.store.IsExist(ctx, key)
}

// Invalidate 通过tag批量删除
func (c *Cache) ClearByTags(ctx context.Context, tags []string, opts ...store.Option) error {
	return c.store.ClearByTags(ctx, tags, opts...)
}

// Invalidate 通过tag批量删除
func (c *Cache) HSet(ctx context.Context, key string, val map[string]interface{}, opts ...store.Option) error {
	return c.store.HSet(ctx, key, val, opts...)
}

// HGet 通过key和map的key获取map指定字段值
func (c *Cache) HGet(ctx context.Context, key, field string) (interface{}, error) {
	return c.store.HGet(ctx, key, field)
}

// HGetAll 返回所有map字段和对应的值
// HGet 通过key和map的key获取map指定字段值
func (c *Cache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.store.HGetAll(ctx, key)
}

// HIncrBy map的field对应value加上一个加数
func (c *Cache) HIncrBy(ctx context.Context, key, field string, incr int64) error {
	return c.store.HIncrBy(ctx, key, field, incr)
}

// Expire 设置过期时间
func (c *Cache) Expire(ctx context.Context, key string, duration time.Duration) error {
	return c.store.Expire(ctx, key, duration)
}
