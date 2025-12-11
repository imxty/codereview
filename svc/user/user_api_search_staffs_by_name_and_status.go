package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 根据员工名和状态模糊查询
func (u *UserAPIHandler) SearchStaffsByNameAndStatus(ctx context.Context, req *pb.SearchStaffsByNameAndStatusRequest, rsp *pb.SearchStaffsByNameAndStatusResponse) error {
	err := validateSearchStaffsByNameAndStatusRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询用户数据
	users, err := u.userStore.SearchStaffsByNameAndStatus(ctx, req.GetTenantId(), req.GetStaffName(), req.GetIsActivated())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回响应
	staffs := make([]*pb.Staff, len(users))
	for i, each := range users {
		staffs[i] = toProtoStaffFromUser(each)
	}
	rsp.Staffs = staffs

	return nil
}

// 验证 request
func validateSearchStaffsByNameAndStatusRequest(req *pb.SearchStaffsByNameAndStatusRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
