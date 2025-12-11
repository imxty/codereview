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

func (s *TenantAPIHandler) UpdateStaff(ctx context.Context, req *pb.UpdateStaffRequest, rsp *pb.UpdateStaffResponse) error {
	err := validateUpdateStaffRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	if req.GetPlainPassword() != "" {
		// 密码长度必须大于6位
		if len(req.GetPlainPassword()) < 6 {
			return errors.Error(codes.InvalidRequest, "密码长度必须大于6位")
		}
	}
	// 更新员工信息请求
	_, err = s.userAPI.UpdateStaff(ctx, &userpb.UpdateStaffRequest{
		// 员工ID
		StaffId: req.GetStaffId(),
		// 手机号
		Phone: req.GetPhone(),
		// 姓名
		Name: req.GetName(),
		// 新的明文密码
		PlainPassword: req.GetPlainPassword(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	// 如果修改员工密码，强制退出
	if req.GetPlainPassword() != "" {
		// 删除该员工的token,强制退出
		err = s.tokenStore.KickOutUser(ctx, req.GetStaffId())
		if err != nil {
			return errors.Error(api.ErrCreateAccessToken, api.ErrorChineseMsg(api.ErrKickedOut))
		}
	}
	return nil
}

// 验证request
func validateUpdateStaffRequest(req *pb.UpdateStaffRequest) error {
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetName() == "" {
		return gerr.New("staff_name should not be empty")
	}
	return nil
}
