package user

import (
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/filestore"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
)

type UserAPIHandler struct {
	userStore       domain.UserRepository
	notificationAPI notificationpb.NotificationAPIService
	reviewAPI       reviewpb.ReviewAPIService
	customerAPI     customerpb.CustomerAPIService
	reportAPI       reportpb.ReportAPIService
	productAPI      productpb.ProductAPIService
	s3Store         filestore.FileStore
}

func (svc *UserAPIHandler) Name() string {
	const name = "UserAPI"
	return name
}

var _ pb.UserAPIHandler = (*UserAPIHandler)(nil)

func NewUserAPIHandler(userStore domain.UserRepository, notificationAPI notificationpb.NotificationAPIService, reviewAPI reviewpb.ReviewAPIService, customerAPI customerpb.CustomerAPIService, reportAPI reportpb.ReportAPIService, productAPI productpb.ProductAPIService, s3Store filestore.FileStore) *UserAPIHandler {
	return &UserAPIHandler{
		userStore:       userStore,
		notificationAPI: notificationAPI,
		reviewAPI:       reviewAPI,
		customerAPI:     customerAPI,
		reportAPI:       reportAPI,
		productAPI:      productAPI,
		s3Store:         s3Store,
	}
}
