package report

import (
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	devicepb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/filestore"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	platformpb "github.com/jinmukeji/platform-proto/gen/go/platform/report/v1"
)

type ReportAPIHandler struct {
	reportStore    domain.ReportRepository
	platformClient platformpb.ReportAPIClient
	customerAPI    customerpb.CustomerAPIService
	deviceAPI      devicepb.DeviceAPIService
	userAPI        userpb.UserAPIService
	productAPI     productpb.ProductAPIService
	appId          string
	reportLink     string

	s3Domain          string
	s3Store           filestore.FileStore
	faceTongueAPI     string
	faceTongueAppCode string
}

var _ pb.ReportAPIHandler = (*ReportAPIHandler)(nil)

func (svc *ReportAPIHandler) Name() string {
	const name = "ReportAPI"
	return name
}

func NewReportAPIHandler(reportStore domain.ReportRepository,
	platformClient platformpb.ReportAPIClient,
	customerAPI customerpb.CustomerAPIService,
	deviceAPI devicepb.DeviceAPIService,
	userAPI userpb.UserAPIService,
	productAPI productpb.ProductAPIService,
	appId string,
	reportLink string,
	s3Store filestore.FileStore,
	faceTongueAPI string,
	faceTongueAppCode string,
	s3Domain string) *ReportAPIHandler {
	return &ReportAPIHandler{
		reportStore:       reportStore,
		platformClient:    platformClient,
		customerAPI:       customerAPI,
		deviceAPI:         deviceAPI,
		userAPI:           userAPI,
		productAPI:        productAPI,
		appId:             appId,
		reportLink:        reportLink,
		s3Store:           s3Store,
		faceTongueAPI:     faceTongueAPI,
		faceTongueAppCode: faceTongueAppCode,
		s3Domain:          s3Domain,
	}
}
