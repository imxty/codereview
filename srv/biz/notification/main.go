package main

import (
	mlog "github.com/jinmukeji/go-pkg/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/huimaibao-service/svc/notification/client"
	"github.com/jinmukeji/plat-pkg/v4/rpc"
	"github.com/jinmukeji/plat-pkg/v4/rpc/service"
	"go-micro.dev/v4"
	mc "go-micro.dev/v4/client"
	"go-micro.dev/v4/server"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"

	nstore "github.com/jinmukeji/huimaibao-service/store/notification"
	notificationsvc "github.com/jinmukeji/huimaibao-service/svc/notification"
)

const (
	// ServiceName 是本微服务的名称
	ServiceName = "notification"
	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.shangyikangyou.huimaibao.service"
)

// 短信配置
type smsConfig struct {
	// 阿里云短信key
	AliyunSmsAccessKeyID string `json:"aliyun_sms_access_key_id" yaml:"aliyun_sms_access_key_id"`
	// 阿里云密钥
	AliyunSmsAccessKeySecret string `json:"aliyun_sms_access_key_secret" yaml:"aliyun_sms_access_key_id"`
	// 环境
	Env string `json:"env" yaml:"env"`
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

	svc := service.CreateService(opts)
	err := service.RegisterServer(svc.Server(), serviceRegister(svc))
	die(err)

	// Run the service
	err = svc.Run()
	die(err)
}

func serviceRegister(service micro.Service) service.RegisterServerFunc {
	return func(srv server.Server) error {

		// 设置重试次数
		err := service.Client().Init(
			mc.Retries(0),
		)
		die(err)

		dbConfigKey := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.service.@global",
			"mysql"}

		// 获取数据库的配置信息
		dbConn := &dbutils.DBConnection{}
		err = rpc.YamlConfig().Get(dbConfigKey...).Scan(dbConn)
		if err != nil {
			die(err)
		}
		// 连接数据库
		db, err := gorm.Open(mysql.Open(dbutils.GetDsn(dbConn)), &gorm.Config{})
		if err != nil {
			die(err)
		}
		// Register Notification Store
		ucNotification := nstore.NewNotificationStore(dbutils.NewConnection(db))

		// 初始化阿里云短信服务
		smsConfigKey := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.service.@global",
			"sms",
			"config"}

		smsConf := &smsConfig{}
		err = rpc.YamlConfig().Get(smsConfigKey...).Scan(smsConf)
		if err != nil {
			return err
		}
		smsClient, err := client.NewAliyunSMSClient(smsConf.AliyunSmsAccessKeyID, smsConf.AliyunSmsAccessKeySecret, smsConf.Env)
		if err != nil {
			return err
		}

		smsSignNameKey := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.service.@global",
			"sms",
			"sign_name"}
		signName := ""
		signName = rpc.YamlConfig().Get(smsSignNameKey...).String(signName)

		// Register services
		notiAPI := notificationsvc.NewNotificationAPIHandler(ucNotification, smsClient, signName)

		if err := pb.RegisterNotificationAPIHandler(srv, notiAPI); err != nil {
			return err
		}
		log.Infof("Registered RPC service: %s", notiAPI.Name())

		return nil
	}
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
