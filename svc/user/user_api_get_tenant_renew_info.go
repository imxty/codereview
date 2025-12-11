package user

import (
	"context"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (u *UserAPIHandler) GetTenantRenewInfo(ctx context.Context, req *pb.GetTenantRenewInfoRequest, rsp *pb.GetTenantRenewInfoResponse) error {
	err := validateGetTenantRenewInfoRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商户组织信息
	organization, err := u.userStore.GetOrganizationTenantByTenantId(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if organization == nil {
		return errors.Errorf(ErrOrganizationTenantNotExist, "organization_tenant[%s] not found", req.GetTenantId())
	}

	// 获取组织联系人手机
	user, err := u.userStore.GetOrganizationByOrganizationId(ctx, organization.GetOrganizationID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if user == nil {
		return errors.Errorf(ErrOrganizationNotExist, "organization[%s] not found", organization.GetOrganizationID())
	}

	// 获取商户订阅信息
	st, err := u.userStore.GetLatestTenantSubscriptionTimeline(ctx, organization.GetOrganizationID(), req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	nowTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	// 没有订阅或30天内不过期
	if st == nil {
		rsp.IsExpiredSoon = false
	} else if nowTime.After(st.GetEndTime()) || (st.GetEndTime().Sub(nowTime).Hours()/24) <= 30 {
		rsp.IsExpiredSoon = true
	} else {
		rsp.IsExpiredSoon = false
	}

	if st == nil {
		rsp.EndTime = nil
	} else {
		rsp.EndTime = timestamppb.New(st.GetEndTime())
	}

	rsp.OrganizationContactPhone = user.GetOrganizationContactPhone()

	return nil
}

func validateGetTenantRenewInfoRequest(req *pb.GetTenantRenewInfoRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
