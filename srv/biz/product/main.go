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

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"

	pstore "github.com/jinmukeji/huimaibao-service/store/product"
	productsvc "github.com/jinmukeji/huimaibao-service/svc/product"
)

type s3Config struct {
	BucketName  string `json:"bucket_name"`
	AccessKeyID string `json:"access_key_id"`
	SecretKey   string `json:"secret_key"`
	Region      string `json:"region"`
	KeyPrefix   string `json:"key_prefix"`
}

const (
	// ServiceName 是本微服务的名称
	ServiceName = "product"
	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.shangyikangyou.huimaibao.service"
	// ReviewServiceName 审核服务名称
	ReviewServiceName = "com.shangyikangyou.huimaibao.service.review"
	// UserServiceName 永辉服务名称
	UserServiceName = "com.shangyikangyou.huimaibao.service.user"
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
		// Register Product Store
		ucProduct := pstore.NewProductStore(dbutils.NewConnection(db))

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

		// Init ProductService
		reviewAPI := reviewpb.NewReviewAPIService(ReviewServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Init UserService
		userAPI := userpb.NewUserAPIService(UserServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Register services
		productAPI := productsvc.NewProductAPIHandler(ucProduct, s3Store, reviewAPI, userAPI)

		if err := pb.RegisterProductAPIHandler(srv, productAPI); err != nil {
			return err
		}
		log.Infof("Registered RPC service: %s", productAPI.Name())

		return nil
	}
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
