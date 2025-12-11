package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取员工列表
func (u *UserAPIHandler) ListStaffs(ctx context.Context, req *pb.ListStaffsRequest, rsp *pb.ListStaffsResponse) error {
	// 验证 request
	err := validateListStaffsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询租户
	tenant, err := u.userStore.GetTenantUser(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenant == nil {
		return errors.Error(ErrTenantNotFound, "tenant not found")
	}

	// 获取租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Error(ErrTenantNotFound, "tenant_entity not found")
	}

	// 获取所有员工信息
	staffs, err := u.userStore.ListAllStaffs(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 构造返回的参数
	activateMemberCount := 0
	tenantInfo := toProtoTenant(tenantEntity)
	pStaffs := make([]*pb.Staff, len(staffs))
	for k, v := range staffs {
		// 属于该租户才去给他赋值当前租户的信息，如果被删除了，那么 is_activated 就是 0
		pStaffs[k] = toProtoStaffFromUser(v)
		if v.GetTenantID() == req.GetTenantId() && v.GetIsActivated() {
			activateMemberCount++
			pStaffs[k].Tenant = tenantInfo
		}
		pStaffs[k].IsActivated = v.GetIsActivated()
		pStaffs[k].IsDeleted = !v.GetIsActivated()
	}
	// 返回最大人数
	rsp.MaxStaffCount = tenantEntity.GetStaffCountQuota()
	rsp.ActivatedStaffCount = int32(activateMemberCount)
	rsp.Staffs = pStaffs

	return nil
}

// 验证 request
func validateListStaffsRequest(req *pb.ListStaffsRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	return nil
}
