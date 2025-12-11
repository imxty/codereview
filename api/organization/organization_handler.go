package organization

import (
	"github.com/go-redis/redis/v8"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api/token"
)

type OrganizationAPIHandler struct {
	userAPI         userpb.UserAPIService
	reviewAPI       reviewpb.ReviewAPIService
	productAPI      productpb.ProductAPIService
	notificationAPI notificationpb.NotificationAPIService
	customerAPI     customerpb.CustomerAPIService
	reportAPI       reportpb.ReportAPIService
	tokenStore      token.TokenStore
	s3Domain        string
	// redis
	redisCli *redis.Client
}

var _ pb.UserAPIHandler = (*OrganizationAPIHandler)(nil)
var _ pb.CustomerAPIHandler = (*OrganizationAPIHandler)(nil)
var _ pb.ReportAPIHandler = (*OrganizationAPIHandler)(nil)
var _ pb.NotificationAPIHandler = (*OrganizationAPIHandler)(nil)
var _ pb.ProductAPIHandler = (*OrganizationAPIHandler)(nil)

func NewOrganizationAPIHandler(user userpb.UserAPIService, customer customerpb.CustomerAPIService,
	reviewAPI reviewpb.ReviewAPIService,
	report reportpb.ReportAPIService,
	notification notificationpb.NotificationAPIService,
	product productpb.ProductAPIService, tokenStore token.TokenStore, s3Domain string, redisCli *redis.Client) *OrganizationAPIHandler {
	return &OrganizationAPIHandler{
		userAPI:         user,
		customerAPI:     customer,
		reviewAPI:       reviewAPI,
		reportAPI:       report,
		notificationAPI: notification,
		productAPI:      product,
		tokenStore:      tokenStore,
		s3Domain:        s3Domain,
		redisCli:        redisCli,
	}
}
