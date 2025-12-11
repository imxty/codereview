package store

import (
	"time"
)

// Option 连接配置
type Option func(*Options)

// Options represents the cache store available options
type Options struct {

	// Expiration allows to specify an expiration time when setting a value
	Expiration time.Duration

	// Tags allows to specify associated tags to the current value
	Tags []string

	// Tag前缀
	TagPattern string
}

// WithExpiration 设置过期时长
func WithExpiration(expiration time.Duration) Option {
	return func(h *Options) {
		h.Expiration = expiration
	}
}

// WithTagPattern 设置tag前缀
func WithTagPattern(pattern string) Option {
	return func(h *Options) {
		h.TagPattern = pattern
	}
}

// WithTags 设置tags
func WithTags(tags []string) Option {
	return func(h *Options) {
		h.Tags = tags
	}
}

const (
	expiration = time.Hour
	parrten    = "redis_tag_%s"
)

// NewOptions 初始化连接选项
func NewOptions(opts ...Option) *Options {
	// 设置默认值
	c := &Options{
		Expiration: expiration,
		Tags:       []string{},
		TagPattern: parrten,
	}

	// 设置自定义配置参数
	c.applyOption(opts...)
	return c
}

// 使用Option
func (o *Options) applyOption(opts ...Option) {
	// 设置自定义配置参数
	for _, opt := range opts {
		opt(o)
	}
}

// ExpirationValue returns the expiration option value
func (o Options) ExpirationValue() time.Duration {
	return o.Expiration
}

// TagsValue returns the tags option value
func (o Options) TagsValue() []string {
	return o.Tags
}

// Pattern returns the tag pattern
func (o Options) Pattern() string {
	return o.TagPattern
}

// deepCopyOptions 深拷贝选项
func (s *Options) deepCopyOptions() *Options {
	var t []string
	// 拷贝一份slice
	copy(t, s.Tags)
	v := &Options{
		Expiration: s.Expiration,
		Tags:       t,
		TagPattern: s.TagPattern,
	}
	return v
}
