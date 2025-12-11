package app

import (
	"context"
	gerr "errors"
	"unicode/utf8"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// 有效手机号码长度
	validPhoneLength = 11
)

func (s *AppAPIHandler) AddCustomer(ctx context.Context, req *pb.AddCustomerRequest, rsp *pb.AddCustomerResponse) error {

	// 验证request
	err := validateAddCustomerRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 数据转化为Svc的数据
	svcCustomer, err := toSvcCustomer(req.GetCustomer())
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}
	// 发送请求
	addCustomerRsp, err := s.customerAPI.AddCustomer(ctx, &customerpb.AddCustomerRequest{
		TenantId: req.GetTenantId(),
		Customer: svcCustomer,
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据给app
	appCustomer, err := toAppCustomer(addCustomerRsp.GetCustomer())
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}

	rsp.Customer = appCustomer

	return nil
}

// 验证request
func validateAddCustomerRequest(req *pb.AddCustomerRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomer().GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetCustomer().GetNickname() == "" || utf8.RuneCountInString(req.GetCustomer().GetNickname()) > 50 {
		return gerr.New("invalid nickname")
	}
	if req.GetCustomer().GetAreaCode() == "" {
		return gerr.New("area_code should not be empty")
	}
	if req.GetCustomer().GetPhone() == "" || len(req.GetCustomer().GetPhone()) != validPhoneLength {
		return gerr.New("invalid phone number")
	}
	if req.GetCustomer().GetBirthday() == nil {
		return gerr.New("birthday should not be nil")
	}
	if req.GetCustomer().GetHeight() != 0 && (req.GetCustomer().GetHeight() < 50 || req.GetCustomer().GetHeight() > 250) {
		return gerr.New("invalid height")
	}
	if req.GetCustomer().GetWeight() != 0 && (req.GetCustomer().GetWeight() < 30 || req.GetCustomer().GetWeight() > 250) {
		return gerr.New("invalid weight")
	}
	if utf8.RuneCountInString(req.GetCustomer().GetPmh()) > 100 {
		return gerr.New("invalid pmh")
	}
	if utf8.RuneCountInString(req.GetCustomer().GetRemarks()) > 100 {
		return gerr.New("invalid remarks")
	}
	return nil
}
