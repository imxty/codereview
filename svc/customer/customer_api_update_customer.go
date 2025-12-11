package customer

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// - nickname -- 常客昵称
// - gender -- 常客性别
// - phone -- 常客手机号
// - birthday -- 常客生日
// - height -- 常客身高
// - weight -- 常客体重
// - pmh -- 常客既往病史
// - remark -- 常客其他备注
func (s *CustomerAPIHandler) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest, rsp *pb.UpdateCustomerResponse) error {

	// 验证request
	err := validateUpdateCustomerRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 构建需要更新的字段和值
	// 包括版本号
	updates, err := generateCustomerUpdates(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获得当前常客版本号
	rCustomer, err := s.customerStore.GetCustomer(ctx, req.GetTenantId(), req.GetCustomer().GetCustomerId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if rCustomer == nil {
		return errors.Errorf(ErrCustomerNotFound, "customer[%s] not found", req.GetCustomer().GetCustomerId())
	}
	// 更新常客
	err = s.customerStore.UpdateCustomer(ctx, req.GetTenantId(), req.GetCustomer().GetCustomerId(), rCustomer.GetRev(), updates)
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "update customer failed[%s]", err.Error())
	}

	return nil
}

// 验证request
func validateUpdateCustomerRequest(req *pb.UpdateCustomerRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomer().GetCustomerId() == "" {
		return gerr.New("customer_id should not be empty")
	}
	return nil
}

// 构建更新数据
func generateCustomerUpdates(req *pb.UpdateCustomerRequest) (map[string]interface{}, error) {
	// 初始化更新条件
	updates := make(map[string]interface{})
	birth := req.GetCustomer().GetBirthday()
	cusAge := getAge(birth.GetYear(), birth.GetMonth(), birth.GetDay())
	if cusAge < 16 {
		return nil, errors.Errorf(ErrCustomerAge, "invalid age[%d]", cusAge)
	}
	var gender int32
	if req.GetCustomer().GetGender() == pb.Gender_GENDER_FEMALE {
		gender = 3
	} else if req.GetCustomer().GetGender() == pb.Gender_GENDER_MALE {
		gender = 2
	} else {
		return nil, gerr.New("invalid gender")
	}
	updates["nickname"] = req.GetCustomer().GetNickname()
	updates["gender"] = gender
	updates["phone"] = req.GetCustomer().GetPhone()
	updates["birthday"] = toDomainBirthday(req.Customer.GetBirthday())
	updates["height"] = req.GetCustomer().GetHeight()
	updates["weight"] = req.GetCustomer().GetWeight()
	updates["pmh"] = req.GetCustomer().GetPmh()
	updates["remark"] = req.GetCustomer().GetRemarks()
	updates["initial"] = getNicknameInitial(req.GetCustomer().GetNickname())

	return updates, nil
}
