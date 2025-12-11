package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
)

func (s *AppAPIHandler) GetStaff(ctx context.Context, req *pb.GetStaffRequest, rsp *pb.GetStaffResponse) error {

	// 1.验证request
	err := validateGetStaffRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 发送获取员工信息的请求
	getStaffRsp, err := s.userAPI.GetStaff(ctx, &userpb.GetStaffRequest{
		StaffId:  req.GetStaffId(),
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 如果被删除返回错误
	if getStaffRsp.GetStaff().GetIsDeleted() {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	rsp.Staff = toAppStaff(getStaffRsp.GetStaff(), s.s3Domain)

	return nil
}

// 验证request
func validateGetStaffRequest(req *pb.GetStaffRequest) error {
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
