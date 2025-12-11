package boss

import (
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api/token"
)

type BossAPIHandler struct {
	userAPI         userpb.UserAPIService
	reviewAPI       reviewpb.ReviewAPIService
	reportAPI       reportpb.ReportAPIService
	notificationAPI notificationpb.NotificationAPIService
	tokenStore      token.TokenStore
	s3Domain        string
}

var _ pb.UserAPIHandler = (*BossAPIHandler)(nil)
var _ pb.ReviewAPIHandler = (*BossAPIHandler)(nil)
var _ pb.ReportAPIHandler = (*BossAPIHandler)(nil)
var _ pb.NotificationAPIHandler = (*BossAPIHandler)(nil)

func NewBossAPIHandler(userAPI userpb.UserAPIService,
	reviewAPI reviewpb.ReviewAPIService,
	reportAPI reportpb.ReportAPIService,
	notificationAPI notificationpb.NotificationAPIService,
	tokenStore token.TokenStore,
	s3Domain string) *BossAPIHandler {
	return &BossAPIHandler{
		userAPI:         userAPI,
		reviewAPI:       reviewAPI,
		reportAPI:       reportAPI,
		notificationAPI: notificationAPI,
		tokenStore:      tokenStore,
		s3Domain:        s3Domain,
	}
}
