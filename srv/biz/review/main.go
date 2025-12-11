package main

import (
	mlog "github.com/jinmukeji/go-pkg/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/rpc"
	"github.com/jinmukeji/plat-pkg/v4/rpc/service"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/server"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"

	rstore "github.com/jinmukeji/huimaibao-service/store/review"
	reviewsvc "github.com/jinmukeji/huimaibao-service/svc/review"
)

const (
	// ServiceName 是本微服务的名称
	ServiceName = "review"
	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.shangyikangyou.huimaibao.service"
	// ProductServiceName 商品服务名称
	ProductServiceName = "com.shangyikangyou.huimaibao.service.product"
	// UserServiceName 永辉服务名称
	UserServiceName = "com.shangyikangyou.huimaibao.service.user"
	// NotificationServiceName 通知服务名称
	NotificationServiceName = "com.shangyikangyou.huimaibao.service.notification"
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

		// ucReview
		ucReview := rstore.NewReviewStore(dbutils.NewConnection(db))

		// Init ProductService
		productAPI := productpb.NewProductAPIService(ProductServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Init UserService
		userAPI := userpb.NewUserAPIService(UserServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Init NotificationService
		notificationAPI := notificationpb.NewNotificationAPIService(NotificationServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Register services
		reviewAPI := reviewsvc.NewReviewAPIHandler(ucReview, productAPI, userAPI, notificationAPI)

		if err := pb.RegisterReviewAPIHandler(srv, reviewAPI); err != nil {
			return err
		}
		log.Infof("Registered RPC service: %s", reviewAPI.Name())

		return nil
	}
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
