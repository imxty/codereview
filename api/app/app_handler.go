package app

import (
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	devicepb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api/appupdate/android"
	"github.com/jinmukeji/huimaibao-service/api/token"
)

type AppAPIHandler struct {
	userAPI         userpb.UserAPIService
	deviceAPI       devicepb.DeviceAPIService
	customerAPI     customerpb.CustomerAPIService
	reportAPI       reportpb.ReportAPIService
	notificationAPI notificationpb.NotificationAPIService
	tokenStore      token.TokenStore
	apkClient       android.ApkS3Client
	s3Domain        string
}

var _ pb.DeviceAPIHandler = (*AppAPIHandler)(nil)
var _ pb.UserAPIHandler = (*AppAPIHandler)(nil)
var _ pb.ReportAPIHandler = (*AppAPIHandler)(nil)
var _ pb.NotificationAPIHandler = (*AppAPIHandler)(nil)
var _ pb.CustomerAPIHandler = (*AppAPIHandler)(nil)

func NewDeskAppHandler(user userpb.UserAPIService, customer customerpb.CustomerAPIService,
	device devicepb.DeviceAPIService,
	reportAPI reportpb.ReportAPIService,
	notificationAPI notificationpb.NotificationAPIService, tokenStore token.TokenStore, s3Domain string, apkClient android.ApkS3Client) *AppAPIHandler {
	return &AppAPIHandler{
		userAPI:         user,
		deviceAPI:       device,
		customerAPI:     customer,
		reportAPI:       reportAPI,
		notificationAPI: notificationAPI,
		tokenStore:      tokenStore,
		s3Domain:        s3Domain,
		apkClient:       apkClient,
	}
}
