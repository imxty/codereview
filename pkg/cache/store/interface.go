package store

import (
	"context"
	"time"
)

// StoreInterface 主要面向的业务包括: token，用户的信息，商品等
type StoreInterface interface {
	// Get 通过key获取内容，找不到不会报错，返回nil
	Get(ctx context.Context, key string) (interface{}, error)
	// GetMulti 获取多个key的内容，找不到不会报错，返回nil
	GetMulti(ctx context.Context, keys []string) ([]interface{}, error)
	// Put 设置key，value
	// 更新就是直接重新put即可
	Put(ctx context.Context, key string, val interface{}, opts ...Option) error
	// Incr value加一
	Incr(ctx context.Context, key string) error
	// Decr value减一
	Decr(ctx context.Context, key string) error
	// Delete 删除key
	Delete(ctx context.Context, key string) error
	// IsExist 检测key是否存在
	IsExist(ctx context.Context, key string) (bool, error)
	// ClearByTags 通过tag批量删除
	ClearByTags(ctx context.Context, tags []string, opts ...Option) error
	// HSet 设置map
	HSet(ctx context.Context, key string, val map[string]interface{}, opts ...Option) error
	// HGet 通过key和map的key获取map指定字段值
	HGet(ctx context.Context, key, field string) (interface{}, error)
	// HGetAll 返回所有map字段和对应的值
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	// HIncrBy map的field对应value加上一个加数
	HIncrBy(ctx context.Context, key, field string, incr int64) error
	// Expire 设置过期时间
	Expire(ctx context.Context, key string, duration time.Duration) error
}
