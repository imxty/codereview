package s3

// Options 是 aws 连接配置
type Options struct {
	// aws 存储桶名称
	BucketName string `json:"bucket_name" yaml:"bucket_name"`
	// aws api 访问身份认证
	AccessKeyID string `json:"access_key" yaml:"access_key"`
	// aws api 访问密钥
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	// aws 所在区域
	Region string `json:"region" yaml:"region"`
	// aws 上数据的存储路径
	KeyPrefix string `json:"key_prefix" yaml:"key_prefix"`
}

// Option 是 aws 连接配置设置方法
type Option func(opt *Options)

// BucketName 设置 aws 存储桶地址
func BucketName(name string) Option {
	return func(options *Options) {
		options.BucketName = name
	}
}

// AccessKeyID 设置 aws api 访问身份认证信息
func AccessKeyID(id string) Option {
	return func(options *Options) {
		options.AccessKeyID = id
	}
}

// SecretKey 设置 aws api 访问密钥
func SecretKey(key string) Option {
	return func(options *Options) {
		options.SecretKey = key
	}
}

// Region 设置 aws 区域
func Region(region string) Option {
	return func(options *Options) {
		options.Region = region
	}
}

// KeyPrefix 设置数据存储路径
func KeyPrefix(keyPrefix string) Option {
	return func(options *Options) {
		options.KeyPrefix = keyPrefix
	}
}

// defaultOptions 返回 aws 默认连接配置
func defaultOptions() *Options {
	return &Options{
		BucketName:  "",
		AccessKeyID: "",
		SecretKey:   "",
		Region:      "cn-north-1",
		KeyPrefix:   "",
	}
}

// newOptions 返回新的 aws 连接配置
func newOptions(opts ...Option) *Options {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}
