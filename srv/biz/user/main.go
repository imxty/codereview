package main

import (
	"encoding/json"

	mlog "github.com/jinmukeji/go-pkg/v2/log"
	s3store "github.com/jinmukeji/huimaibao-service/pkg/filestore/s3"
	"go-micro.dev/v4/client"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/rpc"
	"github.com/jinmukeji/plat-pkg/v4/rpc/service"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"

	ustore "github.com/jinmukeji/huimaibao-service/store/user"
	usersvc "github.com/jinmukeji/huimaibao-service/svc/user"
)

const (
	// ServiceName 是本微服务的名称
	ServiceName = "user"
	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.shangyikangyou.huimaibao.service"
	// NotificationServiceName 短信服务名称
	NotificationServiceName = "com.shangyikangyou.huimaibao.service.notification"
	// ReviewServiceName 审核服务名称
	ReviewServiceName = "com.shangyikangyou.huimaibao.service.review"
	// CustomerServiceName 常客服务名称
	CustomerServiceName = "com.shangyikangyou.huimaibao.service.customer"
	// ReportServiceName 报告服务名称
	ReportServiceName = "com.shangyikangyou.huimaibao.service.report"
	// ProductServiceName 商品服务名称
	ProductServiceName = "com.shangyikangyou.huimaibao.service.product"
)

type s3Config struct {
	BucketName  string `json:"bucket_name" yaml:"bucket_name"`
	AccessKeyID string `json:"access_key_id" yaml:"access_key_id"`
	SecretKey   string `json:"secret_key" yaml:"secret_key"`
	Region      string `json:"region" yaml:"region"`
	KeyPrefix   string `json:"key_prefix" yaml:"key_prefix"`
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
		ucUser := ustore.NewUserStore(dbutils.NewConnection(db))

		fileConfigKey := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.service.@global",
			"file_store"}

		var fileStore s3Config
		err = json.Unmarshal(rpc.YamlConfig().Get(fileConfigKey...).Bytes(), &fileStore)
		if err != nil {
			return err
		}
		// Init S3Service
		s3Store, err := s3store.NewS3Store(
			s3store.BucketName(fileStore.BucketName),
			s3store.AccessKeyID(fileStore.AccessKeyID),
			s3store.SecretKey(fileStore.SecretKey),
			s3store.Region(fileStore.Region),
			s3store.KeyPrefix(fileStore.KeyPrefix),
		)
		if err != nil {
			die(err)
		}

		// Init NotificationService
		notiAPI := notificationpb.NewNotificationAPIService(NotificationServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Init CustomerService
		customerAPI := customerpb.NewCustomerAPIService(CustomerServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Init ReportService
		reportAPI := reportpb.NewReportAPIService(ReportServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Init ReviewService
		reviewAPI := reviewpb.NewReviewAPIService(ReviewServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Init ProductService
		productAPI := productpb.NewProductAPIService(ProductServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Register services
		userAPI := usersvc.NewUserAPIHandler(ucUser, notiAPI, reviewAPI, customerAPI, reportAPI, productAPI, s3Store)

		if err := userpb.RegisterUserAPIHandler(srv, userAPI); err != nil {
			return err
		}
		log.Infof("Registered RPC service: %s", userAPI.Name())

		return nil
	}
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
