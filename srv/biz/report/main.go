package main

import (
	"crypto/tls"
	"encoding/json"

	mlog "github.com/jinmukeji/go-pkg/v2/log"
	"go-micro.dev/v4/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	s3store "github.com/jinmukeji/huimaibao-service/pkg/filestore/s3"
	"github.com/jinmukeji/plat-pkg/v4/rpc"
	"github.com/jinmukeji/plat-pkg/v4/rpc/service"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	devicepb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	platformpb "github.com/jinmukeji/platform-proto/gen/go/platform/report/v1"

	rstore "github.com/jinmukeji/huimaibao-service/store/report"
	reportsvc "github.com/jinmukeji/huimaibao-service/svc/report"
)

const (
	// ServiceName 是本微服务的名称
	ServiceName = "report"
	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.shangyikangyou.huimaibao.service"
	// ProductServiceName 商品服务名称
	ProductServiceName = "com.shangyikangyou.huimaibao.service.product"
	// UserServiceName 用户服务名称
	UserServiceName = "com.shangyikangyou.huimaibao.service.user"
	// CustomerServiceName 常客服务名称
	CustomerServiceName = "com.shangyikangyou.huimaibao.service.customer"
	// DeviceServiceName 设备服务名称
	DeviceServiceName = "com.shangyikangyou.huimaibao.service.device"
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

// 平台配置
type platformOptions struct {
	PlatformAddr string `json:"platform_addr"`
}

type s3Config struct {
	BucketName  string `json:"bucket_name" yaml:"bucket_name"`
	AccessKeyID string `json:"access_key_id" yaml:"access_key_id"`
	SecretKey   string `json:"secret_key" yaml:"secret_key"`
	Region      string `json:"region" yaml:"region"`
	KeyPrefix   string `json:"key_prefix" yaml:"key_prefix"`
}

type AliFaceTongueConfig struct {
	AppCode       string `json:"app_code" yaml:"app_code"`
	FaceTongueAPI string `json:"face_tongue_api" yaml:"face_tongue_api"`
}

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

		// ucReport
		ucReport := rstore.NewReportStore(dbutils.NewConnection(db))

		// 读取平台对于huimaibao的appId
		huimaibaoAppId := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.service.@global",
			"report",
			"app_id"}
		var appId string
		// 获取平台配置信息
		appId = rpc.YamlConfig().Get(huimaibaoAppId...).String(appId)

		// 获取分享报告链接
		huimaibaoReportLink := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.service.@global",
			"report",
			"report_share_link"}
		var rl string
		// 获取平台配置信息
		rl = rpc.YamlConfig().Get(huimaibaoReportLink...).String(rl)

		// 读取平台服务器的地址
		platformConfigKey := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.service.@global",
			"report",
			"platform"}
		platformOpt := &platformOptions{}
		// 获取calc配置信息
		err = rpc.YamlConfig().Get(platformConfigKey...).Scan(platformOpt)
		if err != nil {
			return err
		}

		var platConn *grpc.ClientConn

		// 如果为空说明没有密钥直接用Insecure即可
		platConn, err = grpc.Dial(platformOpt.PlatformAddr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		if err != nil {
			log.Panicf("Failed to connect platform. Error: %v", err)
		}
		platformClient := platformpb.NewReportAPIClient(platConn)

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

		ftConfigKey := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.service.@global",
			"ali_ft_config"}

		var ftConfig AliFaceTongueConfig
		err = json.Unmarshal(rpc.YamlConfig().Get(ftConfigKey...).Bytes(), &ftConfig)
		if err != nil {
			return err
		}

		// 获取s3Domain的配置
		s3DomainConfigKey := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.api.@global",
			"s3Domain"}
		var s3Domain string
		s3Domain = rpc.YamlConfig().Get(s3DomainConfigKey...).String(s3Domain)

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

		// Init CustomerService
		customerAPI := customerpb.NewCustomerAPIService(CustomerServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Init DeviceService
		deviceAPI := devicepb.NewDeviceAPIService(DeviceServiceName, service.Client())
		if err != nil {
			die(err)
		}

		// Register services
		reportAPI := reportsvc.NewReportAPIHandler(ucReport, platformClient, customerAPI, deviceAPI, userAPI, productAPI, appId, rl, s3Store, ftConfig.FaceTongueAPI, ftConfig.AppCode, s3Domain)

		if err := pb.RegisterReportAPIHandler(srv, reportAPI); err != nil {
			return err
		}
		log.Infof("Registered RPC service: %s", reportAPI.Name())

		return nil
	}
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
