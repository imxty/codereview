package marshaler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jinmukeji/huimaibao-service/pkg/cache"
	"github.com/jinmukeji/huimaibao-service/pkg/cache/store"
)

// Marshaler is the struct that marshal and unmarshal cache values
type Marshaler struct {
	cache cache.CacheInterface
}

// New creates a new marshaler that marshals/unmarshals cache values
func New(cache cache.CacheInterface) *Marshaler {
	return &Marshaler{
		cache: cache,
	}
}

// Get obtains a value from cache and unmarshal value with given object
func (c *Marshaler) Get(ctx context.Context, key string, returnObj interface{}) error {
	result, err := c.cache.Get(ctx, key)
	if err != nil {
		return err
	}

	return c.unmarshal(result, returnObj)
}

// Put sets a value in cache by marshaling value
func (c *Marshaler) Put(ctx context.Context, key string, object interface{}, opts ...store.Option) error {
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return c.cache.Put(ctx, key, bytes, opts...)
}

// GetMulti 获取多个key的内容，找不到不会报错，返回nil
func (c *Marshaler) GetMulti(ctx context.Context, keys []string, returnObj interface{}) error {
	result, err := c.cache.GetMulti(ctx, keys)
	if err != nil {
		return err
	}
	return c.unmarshal(result, returnObj)
}

// Sets field in the hash stored at key to value.
func (c *Marshaler) HSet(ctx context.Context, key string, object map[string]interface{}, opts ...store.Option) error {
	return c.cache.HSet(ctx, key, object, opts...)
}

// HGet 通过key和map的key获取map指定字段值
func (c *Marshaler) HGet(ctx context.Context, key, field string) (interface{}, error) {
	return c.cache.HGet(ctx, key, field)
}

// Delete 删除key
func (c *Marshaler) Delete(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, key)
}

// IsExist 检测key是否存在
func (c *Marshaler) IsExist(ctx context.Context, key string) (bool, error) {
	return c.cache.IsExist(ctx, key)
}

// ClearByTags 通过tag批量删除
func (c *Marshaler) ClearByTags(ctx context.Context, tags []string, opts ...store.Option) error {
	return c.cache.ClearByTags(ctx, tags, opts...)
}

// 反序列化
func (c *Marshaler) unmarshal(source, target interface{}) error {
	var err error
	switch v := source.(type) {
	case []byte:
		err = json.Unmarshal(v, target)
	case string:
		err = json.Unmarshal([]byte(v), target)
	}

	if err != nil {
		return err
	}
	return nil
}

// HGetAll returns all fields and values of the hash stored at key.
func (c *Marshaler) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.cache.HGetAll(ctx, key)
}

// HIncrBy Increments the number stored at field in the hash stored at key by increment.
func (c *Marshaler) HIncrBy(ctx context.Context, key, field string, incr int64) error {
	return c.cache.HIncrBy(ctx, key, field, incr)
}

// Expire 设置过期时间
func (c *Marshaler) Expire(ctx context.Context, key string, duration time.Duration) error {
	return c.cache.Expire(ctx, key, duration)
}
