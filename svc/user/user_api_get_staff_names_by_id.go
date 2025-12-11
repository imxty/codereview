package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) GetStaffNamesByID(ctx context.Context, req *pb.GetStaffNamesByIDRequest, rsp *pb.GetStaffNamesByIDResponse) error {
	err := validateGetStaffNamesByIDRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	if len(req.GetStaffIds()) < 1 {
		rsp.StaffNames = make(map[string]string)
		return nil
	}

	staffs, err := u.userStore.BatchGetStaffs(ctx, req.GetStaffIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	results := make(map[string]string)
	for _, each := range staffs {
		results[each.GetUserID()] = each.GetNickname()
	}

	rsp.StaffNames = results
	return nil
}

func validateGetStaffNamesByIDRequest(req *pb.GetStaffNamesByIDRequest) error {
	if req.GetStaffIds() == nil {
		return gerr.New("staff_ids should not be nil")
	}
	return nil
}
