package tenant

import (
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"

	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api/token"
)

type TenantAPIHandler struct {
	reportAPI       reportpb.ReportAPIService
	customerAPI     customerpb.CustomerAPIService
	userAPI         userpb.UserAPIService
	reviewAPI       reviewpb.ReviewAPIService
	notificationAPI notificationpb.NotificationAPIService
	tokenStore      token.TokenStore
	s3Domain        string
}

var _ pb.CustomerAPIHandler = (*TenantAPIHandler)(nil)
var _ pb.NotificationAPIHandler = (*TenantAPIHandler)(nil)
var _ pb.ReportAPIHandler = (*TenantAPIHandler)(nil)
var _ pb.TenantAPIHandler = (*TenantAPIHandler)(nil)

func (svc *TenantAPIHandler) Name() string {
	const name = "TenantAPI"
	return name
}

func NewTenantAPIHandler(reportAPI reportpb.ReportAPIService,
	customerAPI customerpb.CustomerAPIService,
	userAPI userpb.UserAPIService,
	reviewAPI reviewpb.ReviewAPIService,
	notificationAPI notificationpb.NotificationAPIService,
	tokenStore token.TokenStore,
	s3Domain string) *TenantAPIHandler {
	return &TenantAPIHandler{
		reportAPI:       reportAPI,
		customerAPI:     customerAPI,
		userAPI:         userAPI,
		reviewAPI:       reviewAPI,
		notificationAPI: notificationAPI,
		tokenStore:      tokenStore,
		s3Domain:        s3Domain,
	}
}
