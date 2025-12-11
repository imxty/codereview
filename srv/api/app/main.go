package main

import (
	"crypto/rsa"
	"encoding/json"
	"time"

	mlog "github.com/jinmukeji/go-pkg/v2/log"
	"github.com/jinmukeji/huimaibao-service/api/app"
	"github.com/jinmukeji/huimaibao-service/api/appupdate/android"
	"go-micro.dev/v4/client"

	deskpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	devicepb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"

	"github.com/jinmukeji/huimaibao-service/api/token"
	"github.com/jinmukeji/huimaibao-service/api/token/tokenstore"
	"github.com/jinmukeji/plat-pkg/v4/auth/jwt"
	"github.com/jinmukeji/plat-pkg/v4/rpc"
	"github.com/jinmukeji/plat-pkg/v4/rpc/service"
	"go-micro.dev/v4"
	server "go-micro.dev/v4/server"
)

const (
	// ServiceName 是本微服务的名称
	ServiceName = "app"
	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.shangyikangyou.huimaibao.api"
	// DeviceServiceName 设备服务名称
	DeviceServiceName = "com.shangyikangyou.huimaibao.service.device"
	// UserServiceName 用户服务名称
	UserServiceName = "com.shangyikangyou.huimaibao.service.user"
	// CustomerServiceName 常客散客服务名称
	CustomerServiceName = "com.shangyikangyou.huimaibao.service.customer"
	// SmsServiceName 短信服务名称
	NotificationServiceName = "com.shangyikangyou.huimaibao.service.notification"
	// ReportServiceName 报告服务名称
	ReportServiceName = "com.shangyikangyou.huimaibao.service.report"
)

// jwt公钥
type jwtOptions struct {
	PublicKey string `json:"public_key"`
}

// redis配置信息
type redisOptions struct {
	RedisDSN string `json:"redis_dsn" yaml:"redis_dsn"`
	DB       int    `json:"db" yaml:"db"`
}

var (
	// jwt公钥内容
	jwtPublicKey *rsa.PublicKey
	// redis 数据源
	redisDSN string
	// redis 库
	redisDB int
	// token 处理
	tokenStore token.TokenStore
	// s3Domain
	s3Domain string
	// 中间件
	wrapper *app.AuthWrapper
	// apkStore
	apkStore *android.ApkOptions
)

// token配置信息
type tokenConfig struct {
	// accessToken 过期时间单位秒
	AtExpireInterval int `json:"at_expire_interval" yaml:"at_expire_interval"`
	// refreshToken 过期时间单位秒
	RtExpireInterval int `json:"rt_expire_interval" yaml:"rt_expire_interval"`
}

var (
	log = mlog.StandardLogger()

	// Following values will be set during build.
	// Do NOT manually modify them.

	// ProductVersion is current product version.
	ProductVersion = "(n/a)"
	// GitCommit is the git commit short hash
	GitCommit = "(n/a)"
	// GoVersion is go compiler version `go version`
	GoVersion = "(n/a)"
	// BuildTime is go build time
	BuildTime = "(n/a)"
)

func main() {
	// ServiceOptions
	opts := service.NewServiceOptions(ServiceNamespace, ServiceName)
	opts.ProductVersion = ProductVersion
	opts.GitCommit = GitCommit
	opts.GoVersion = GoVersion
	opts.BuildTime = BuildTime
	// 加入jwt的检测
	wrapper = app.NewAuthWrapper(&tokenStore, nil)
	opts.PostServerHandlerWrappers = []server.HandlerWrapper{wrapper.Auth}
	svc := service.CreateService(opts)

	// 获取redis的配置
	redisConfigKey := []string{
		"micro",
		"config",
		"jm",
		"com.shangyikangyou.huimaibao.api.@global",
		"redis"}
	redisOpt := &redisOptions{}
	// 获取redis配置信息
	err := rpc.YamlConfig().Get(redisConfigKey...).Scan(redisOpt)
	if err != nil {
		die(err)
	}
	redisDSN = redisOpt.RedisDSN
	redisDB = redisOpt.DB

	// 获取apkConfigKey的配置
	apkConfigKey := []string{
		"micro",
		"config",
		"jm",
		"com.shangyikangyou.huimaibao.api.@global",
		"apk"}
	// 获取apk配置信息
	err = json.Unmarshal(rpc.YamlConfig().Get(apkConfigKey...).Bytes(), &apkStore)
	if err != nil {
		die(err)
	}
	// 读取公钥文件路径
	jwtConfigKey := []string{
		"micro",
		"config",
		"jm",
		"com.shangyikangyou.huimaibao.api.@global",
		"jwt"}
	jwtOpt := &jwtOptions{}
	// 获取jwt公钥
	err = rpc.YamlConfig().Get(jwtConfigKey...).Scan(jwtOpt)
	if err != nil {
		die(err)
	}
	// 读取pk
	jwtPublicKey, err = jwt.LoadRSAPublicKey([]byte(jwtOpt.PublicKey))
	if err != nil {
		die(err)
	}
	wrapper.SetPublicKey(jwtPublicKey)
	wrapper.SetRedisDsn(redisDSN, redisDB)

	// 获取token的时长配置
	tokenConfigKey := []string{
		"micro",
		"config",
		"jm",
		"com.shangyikangyou.huimaibao.api.@global",
		"token"}
	tConfig := &tokenConfig{}
	// 获取redis配置信息
	err = rpc.YamlConfig().Get(tokenConfigKey...).Scan(tConfig)
	if err != nil {
		die(err)
	}
	tokenStore = tokenstore.NewTokenStore(redisDSN, redisDB, time.Second*time.Duration(tConfig.AtExpireInterval), time.Second*time.Duration(tConfig.RtExpireInterval))
	err = service.RegisterServer(svc.Server(), appRegister(svc))
	die(err)

	// Run the service
	err = svc.Run()
	die(err)
}

func appRegister(service micro.Service) service.RegisterServerFunc {
	return func(srv server.Server) error {

		// 设置重试次数
		err := service.Client().Init(
			client.Retries(0),
		)
		die(err)

		// 获取s3Domain的配置
		s3DomainConfigKey := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.api.@global",
			"s3Domain"}
		// 获取redis配置信息
		s3Domain = rpc.YamlConfig().Get(s3DomainConfigKey...).String(s3Domain)

		// NewApkClient
		apkClient, err := android.NewApkClient(apkStore)
		if err != nil {
			return err
		}

		// Init DeviceService
		deviceAPI := devicepb.NewDeviceAPIService(DeviceServiceName, service.Client())

		// Init UserService
		userAPI := userpb.NewUserAPIService(UserServiceName, service.Client())

		// Init CustomerService
		customerAPI := customerpb.NewCustomerAPIService(CustomerServiceName, service.Client())

		// Init reportAPI
		reportAPI := reportpb.NewReportAPIService(ReportServiceName, service.Client())

		// Init notificationAPI
		notificationAPI := notificationpb.NewNotificationAPIService(NotificationServiceName, service.Client())

		deskappAPI := app.NewDeskAppHandler(userAPI, customerAPI, deviceAPI, reportAPI, notificationAPI, tokenStore, s3Domain, apkClient)

		// 注册设备handler
		if err := deskpb.RegisterDeviceAPIHandler(srv, deskappAPI); err != nil {
			return err
		}

		// 注册用户handler
		if err := deskpb.RegisterUserAPIHandler(srv, deskappAPI); err != nil {
			return err
		}

		if err := deskpb.RegisterCustomerAPIHandler(srv, deskappAPI); err != nil {
			return err
		}

		if err := deskpb.RegisterReportAPIHandler(srv, deskappAPI); err != nil {
			return err
		}

		if err := deskpb.RegisterNotificationAPIHandler(srv, deskappAPI); err != nil {
			return err
		}

		return nil
	}
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
