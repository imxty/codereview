package product

import (
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/filestore"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
)

type ProductAPIHandler struct {
	productStore domain.ProductRepository
	s3Store      filestore.FileStore
	reviewAPI    reviewpb.ReviewAPIService
	userAPI      userpb.UserAPIService
}

func (svc *ProductAPIHandler) Name() string {
	const name = "PsroductAPI"
	return name
}

var _ pb.ProductAPIHandler = (*ProductAPIHandler)(nil)

func NewProductAPIHandler(productStore domain.ProductRepository,
	s3Store filestore.FileStore,
	reviewAPI reviewpb.ReviewAPIService,
	userAPI userpb.UserAPIService) *ProductAPIHandler {
	return &ProductAPIHandler{
		productStore: productStore,
		s3Store:      s3Store,
		reviewAPI:    reviewAPI,
		userAPI:      userAPI,
	}
}
