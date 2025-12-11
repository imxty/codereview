package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *BossAPIHandler) UpdateStaffInfo(ctx context.Context, req *pb.UpdateStaffInfoRequest, rsp *pb.UpdateStaffInfoResponse) error {
	err := validateUpdateStaffInfo(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	listRsp, err := s.userAPI.ListSystemUsersByIds(ctx, &userv1.ListSystemUsersByIdsRequest{
		SystemUserIds: []string{req.GetStaffId()},
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	user := listRsp.GetSystemUsers()[req.GetStaffId()]

	_, err = s.userAPI.UpdateStaffInfo(ctx, &userv1.UpdateStaffInfoRequest{
		UserId:      req.GetUserId(),
		StaffId:     req.GetStaffId(),
		StaffName:   req.GetStaffName(),
		StaffPhone:  req.GetStaffPhone(),
		PrivilegeId: req.GetPrivilegeId(),
		Remark:      req.GetRemark(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	if user.GetPrivilegeGroup().GetPrivilegeId() != req.GetPrivilegeId() {
		err := s.tokenStore.KickOutUser(ctx, req.GetStaffId())
		if err != nil {
			return errors.Error(api.ErrCreateAccessToken, api.ErrorChineseMsg(api.ErrKickedOut))
		}
	}

	return nil
}

func validateUpdateStaffInfo(req *pb.UpdateStaffInfoRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetStaffName() == "" {
		return gerr.New("staff_name should not be empty")
	}
	if req.GetStaffPhone() == "" {
		return gerr.New("staff_phone should not be empty")
	}
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
