package review

import (
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
)

type ReviewAPIHandler struct {
	reviewStore     domain.ReviewRepository
	productAPI      productpb.ProductAPIService
	userAPI         userpb.UserAPIService
	notificationAPI notificationpb.NotificationAPIService
}

func (svc *ReviewAPIHandler) Name() string {
	const name = "ReviewAPI"
	return name
}

var _ pb.ReviewAPIHandler = (*ReviewAPIHandler)(nil)

func NewReviewAPIHandler(reviewStore domain.ReviewRepository,
	productAPI productpb.ProductAPIService,
	userAPI userpb.UserAPIService,
	notificationAPI notificationpb.NotificationAPIService) *ReviewAPIHandler {
	return &ReviewAPIHandler{
		reviewStore:     reviewStore,
		productAPI:      productAPI,
		userAPI:         userAPI,
		notificationAPI: notificationAPI,
	}
}
