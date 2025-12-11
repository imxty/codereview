package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisStore redis缓存
type RedisStore struct {
	client *redis.Client
	opt    *Options
}

// NewStore 初始化一个redis的缓存
func NewRedisStore(addr string, db int, opt *Options) (StoreInterface, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
	// 无法连接也会报错
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	// 如果为空，则使用默认的
	finalOpt := NewOptions()
	if opt != nil {
		finalOpt = opt
	}
	return &RedisStore{
		client: client,
		opt:    finalOpt,
	}, nil
}

// 实现 Cache 行为
var _ StoreInterface = (*RedisStore)(nil)

// Put 设置key，value
func (r *RedisStore) Put(ctx context.Context, key string, value interface{}, opts ...Option) error {
	// 深拷贝一份配置
	finalOpt := r.opt.deepCopyOptions()
	finalOpt.applyOption(opts...)
	result := r.client.Set(ctx, key, value, finalOpt.ExpirationValue())

	// 如果有tag，就去为这个key搭上tag
	if tags := finalOpt.TagsValue(); len(tags) > 0 {
		err := r.setTags(ctx, key, finalOpt.Pattern(), tags)
		if err != nil {
			return err
		}
	}
	return result.Err()
}

// Incr value加一
func (r *RedisStore) Incr(ctx context.Context, key string) error {
	return r.client.Incr(ctx, key).Err()
}

// Decr value减一
func (r *RedisStore) Decr(ctx context.Context, key string) error {
	return r.client.Decr(ctx, key).Err()
}

// HIncrBy map的field对应value加上一个加数
func (r *RedisStore) HIncrBy(ctx context.Context, key, field string, incr int64) error {
	return r.client.HIncrBy(ctx, key, field, incr).Err()
}

// GetMulti 获取多个key的内容，找不到不会报错，返回nil
func (r *RedisStore) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	result := r.client.MGet(ctx, keys...)
	if result.Err() != nil && errors.Is(result.Err(), redis.Nil) {
		return nil, nil
	}
	return result.Result()
}

// Get 通过key获取内容，找不到不会报错，返回nil
func (r *RedisStore) Get(ctx context.Context, key string) (interface{}, error) {
	result := r.client.Get(ctx, key)
	if result.Err() != nil && errors.Is(result.Err(), redis.Nil) {
		return nil, nil
	}
	return result.Result()
}

// Delete 删除key
func (r *RedisStore) Delete(ctx context.Context, key string) error {
	result := r.client.Del(ctx, key)
	return result.Err()
}

// IsExist 检测key是否存在
func (r *RedisStore) IsExist(ctx context.Context, key string) (bool, error) {
	result := r.client.Exists(ctx, key)
	if result.Err() != nil {
		return false, result.Err()
	}
	return result.Val() == 1, nil
}

// ClearByTags 通过tags批量删除
func (r *RedisStore) ClearByTags(ctx context.Context, tags []string, opts ...Option) error {
	// 深拷贝一份配置
	finalOpt := r.opt.deepCopyOptions()
	finalOpt.applyOption(opts...)
	var keys []string
	for _, tag := range tags {
		var tagKey = r.buildTag(finalOpt.Pattern(), tag)
		// 找到该tag下所有的key
		result := r.client.SMembers(ctx, tagKey)
		if result.Err() != nil {
			return result.Err()
		}
		// 记录key
		keys = append(keys, result.Val()...)
		keys = append(keys, tagKey)
	}
	// 删除key
	result := r.client.Del(ctx, keys...)

	return result.Err()
}

// setTags 为key设置tags
func (s *RedisStore) setTags(ctx context.Context, key interface{}, tagPattern string, tags []string) error {
	pipe := s.client.TxPipeline()
	for _, tag := range tags {
		var tagKey = s.buildTag(tagPattern, tag)
		pipe.SAdd(ctx, tagKey, key)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// buildTag 构建tag
func (s *RedisStore) buildTag(pattern, tag string) string {
	return fmt.Sprintf(pattern, tag)
}

// HSet 设置map
func (r *RedisStore) HSet(ctx context.Context, key string, val map[string]interface{}, opts ...Option) error {
	if len(val) == 0 {
		return errors.New("value in hset cannot be empty")
	}
	result := r.client.HSet(ctx, key, val)
	return result.Err()
}

// HGet 通过key和map的key获取map指定字段值
func (r *RedisStore) HGet(ctx context.Context, key, field string) (interface{}, error) {
	result := r.client.HGet(ctx, key, field)
	if result.Err() != nil && errors.Is(result.Err(), redis.Nil) {
		return nil, nil
	}
	return result.Result()
}

// HGetAll 返回所有map字段和对应的值
func (r *RedisStore) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result := r.client.HGetAll(ctx, key)
	if result.Err() != nil && errors.Is(result.Err(), redis.Nil) {
		return nil, nil
	}
	return result.Result()
}

// Expire 设置过期时间
func (r *RedisStore) Expire(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Expire(ctx, key, duration).Err()
}
