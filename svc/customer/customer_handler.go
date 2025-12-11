package customer

import (
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/customer/domain"
)

type CustomerAPIHandler struct {
	customerStore domain.CustomerRepository
	userAPI       userpb.UserAPIService
}

func (svc *CustomerAPIHandler) Name() string {
	const name = "CustomerAPI"
	return name
}

var _ pb.CustomerAPIHandler = (*CustomerAPIHandler)(nil)

func NewCustomerAPIHandler(customerStore domain.CustomerRepository, userAPI userpb.UserAPIService) *CustomerAPIHandler {
	return &CustomerAPIHandler{
		customerStore: customerStore,
		userAPI:       userAPI,
	}
}
