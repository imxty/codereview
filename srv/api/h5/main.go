package main

import (
	mlog "github.com/jinmukeji/go-pkg/v2/log"
	h5pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/h5/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"

	"github.com/jinmukeji/huimaibao-service/api/h5"
	"github.com/jinmukeji/plat-pkg/v4/rpc"
	"github.com/jinmukeji/plat-pkg/v4/rpc/service"
	"go-micro.dev/v4"
	server "go-micro.dev/v4/server"
)

const (
	// ServiceName 是本微服务的名称
	ServiceName = "h5"
	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.shangyikangyou.huimaibao.api"
	// ReportServiceName 报告服务名称
	ReportServiceName = "com.shangyikangyou.huimaibao.service.report"
	// ProductServiceName 商品服务名称
	ProductServiceName = "com.shangyikangyou.huimaibao.service.product"
	// UserServiceName 用户服务名称
	UserServiceName = "com.shangyikangyou.huimaibao.service.user"
)

var (
	// s3Domain
	s3Domain string
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

	err := service.RegisterServer(svc.Server(), appRegister(svc))
	die(err)

	// Run the service
	err = svc.Run()
	die(err)
}

func appRegister(service micro.Service) service.RegisterServerFunc {
	return func(srv server.Server) error {

		// 获取s3Domain的配置
		s3DomainConfigKey := []string{
			"micro",
			"config",
			"jm",
			"com.shangyikangyou.huimaibao.api.@global",
			"s3Domain"}
		s3Domain = rpc.YamlConfig().Get(s3DomainConfigKey...).String(s3Domain)

		// Init reportAPI
		reportAPI := reportpb.NewReportAPIService(ReportServiceName, service.Client())

		// Init productAPI
		productAPI := productpb.NewProductAPIService(ProductServiceName, service.Client())

		// Init userAPI
		userAPI := userpb.NewUserAPIService(UserServiceName, service.Client())

		deskappAPI := h5.NewH5APIHandler(reportAPI, productAPI, userAPI, s3Domain)

		if err := h5pb.RegisterReportAPIHandler(srv, deskappAPI); err != nil {
			return err
		}

		if err := h5pb.RegisterProductAPIHandler(srv, deskappAPI); err != nil {
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
