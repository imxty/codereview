package android

// ApkOptions 是 android apk 配置
type ApkOptions struct {
	// aws 存储桶名称
	BucketName string `json:"bucket_name"`
	// aws api 访问身份认证
	AccessKeyID string `json:"access_key_id"`
	// aws api 访问密钥
	SecretKey string `json:"secret_key"`
	// aws 所在区域
	Region string `json:"region"`
	// 安卓apk存放前缀
	ApkPath string `json:"apk_path"`
}

// Apk的相关信息
type ApkInfo struct {
	// 版本号
	Version int32 `json:"version"`
	// 更新信息
	UpdateInfo string `json:"update_info"`
	// apk的大小
	ApkSize string `json:"apk_size"`
	// apk的名字
	ApkName string `json:"apk_name"`
	// 获取apk的链接
	ApkLink string `json:"apk_link"`
}
