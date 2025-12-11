package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) AddStaff(ctx context.Context, req *pb.AddStaffRequest, rsp *pb.AddStaffResponse) error {
	err := validateAddStaffRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 密码长度必须大于6位
	if len(req.GetPlainPassword()) < 6 {
		return errors.Error(codes.InvalidRequest, "密码长度必须大于6位")
	}
	_, err = s.userAPI.AddStaff(ctx, &userpb.AddStaffRequest{
		// 租户ID
		TenantId: req.GetTenantId(),
		// 手机号
		Phone: req.GetPhone(),
		// 员工姓名
		Name: req.GetName(),
		// 初始化密码
		PlainPassword: req.GetPlainPassword(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateAddStaffRequest(req *pb.AddStaffRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetName() == "" {
		return gerr.New("name should not be empty")
	}
	if req.GetPlainPassword() == "" {
		return gerr.New("plain_password should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenantId")
	}
	return nil
}
