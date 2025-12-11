package main

import (
	mlog "github.com/jinmukeji/go-pkg/v2/log"
	"go-micro.dev/v4/client"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/rpc"
	"github.com/jinmukeji/plat-pkg/v4/rpc/service"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"

	dstore "github.com/jinmukeji/huimaibao-service/store/device"
	devicesvc "github.com/jinmukeji/huimaibao-service/svc/device"
)

const (
	// ServiceName 是本微服务的名称
	ServiceName = "device"
	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.shangyikangyou.huimaibao.service"
)

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
			client.Retries(0),
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
		// Register Report Store
		ucDevice := dstore.NewDeviceStore(dbutils.NewConnection(db))

		// Register services
		deviceAPI := devicesvc.NewDeviceAPIHandler(ucDevice)

		if err := pb.RegisterDeviceAPIHandler(srv, deviceAPI); err != nil {
			return err
		}
		log.Infof("Registered RPC service: %s", deviceAPI.Name())

		return nil
	}
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
