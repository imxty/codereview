package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 批量获取员工信息
func (u *UserAPIHandler) BatchGetStaffs(ctx context.Context, req *pb.BatchGetStaffsRequest, rsp *pb.BatchGetStaffsResponse) error {
	// 验证 request
	err := validateBatchGetStaffsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 没有数据返回空
	if len(req.GetStaffIds()) == 0 {
		rsp.Staffs = map[string]*pb.Staff{}
		return nil
	}

	// 批量查询员工
	staffs, err := u.userStore.BatchGetStaffs(ctx, req.GetStaffIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 构建为 map 返回
	staffMap := make(map[string]*pb.Staff)
	for _, v := range staffs {
		staffMap[v.GetUserID()] = &pb.Staff{
			StaffId:     v.GetUserID(),
			Name:        v.GetNickname(),
			Phone:       v.GetPhone(),
			IsActivated: v.GetIsActivated(),
			IsDeleted:   !v.GetIsActivated(),
		}
	}
	rsp.Staffs = staffMap
	return nil
}

// 验证 request
func validateBatchGetStaffsRequest(req *pb.BatchGetStaffsRequest) error {
	if req.GetStaffIds() == nil {
		return gerr.New("staff ids should not be nil")
	}
	return nil
}
