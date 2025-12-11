package main

import (
	mlog "github.com/jinmukeji/go-pkg/v2/log"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	cstore "github.com/jinmukeji/huimaibao-service/store/customer"
	customersvc "github.com/jinmukeji/huimaibao-service/svc/customer"
	"github.com/jinmukeji/plat-pkg/v4/rpc"
	"github.com/jinmukeji/plat-pkg/v4/rpc/service"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/server"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	// ServiceName 是本微服务的名称
	ServiceName = "customer"
	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.shangyikangyou.huimaibao.service"
	// UserServiceName 用户服务名称
	UserServiceName = "com.shangyikangyou.huimaibao.service.user"
	// CustomerServiceName 常客服务名称
	CustomerServiceName = "com.shangyikangyou.huimaibao.service.customer"
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
		ucCustomer := cstore.NewCustomerStore(dbutils.NewConnection(db))

		// Init UserAPIService
		userAPI := userpb.NewUserAPIService(UserServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Register services
		customerAPI := customersvc.NewCustomerAPIHandler(ucCustomer, userAPI)

		if err := customerpb.RegisterCustomerAPIHandler(srv, customerAPI); err != nil {
			return err
		}
		log.Infof("Registered RPC service: %s", customerAPI.Name())

		return nil
	}
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
