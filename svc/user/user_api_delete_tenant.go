package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 删除商户请求
func (u *UserAPIHandler) DeleteTenant(ctx context.Context, req *pb.DeleteTenantRequest, rsp *pb.DeleteTenantResponse) error {
	// 验证请求
	err := validateDeleteTenantRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商户状态
	ot, err := u.userStore.GetOrganizationTenant(ctx, req.GetOrganizationId(), req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if ot == nil {
		return errors.Errorf(ErrOrganizationNotExist, "organization tenant not found %s %s", req.GetOrganizationId(), req.GetTenantId())
	}

	// 如果是资质审核中删除相关 issue
	if ot.GetReviewStatus() == domain.TenantReviewStatusReviewing {
		// 查询 tenant_revision
		revision, err := u.userStore.GetLatestTenantEntityRevision(ctx, req.GetTenantId())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 如果不存在说明商户 id 有问题
		if revision == nil {
			return errors.Errorf(codes.InvalidRequest, "revision not found, can not commit tenant entity[%s]", req.GetTenantId())
		}
		// 取消 issue
		err = u.userStore.CancelRevisionIssue(ctx, revision.GetTenantEntityRevisionID())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	// 1. 删除组织与商户的关系
	// 2. 删除商户角色
	ctx = u.userStore.BeginTx(ctx)
	err = u.userStore.DeleteOrganizationTenant(ctx, req.GetOrganizationId(), req.GetTenantId())
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	err = u.userStore.DeleteTenantUser(ctx, req.GetTenantId())
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	u.userStore.CommitTx(ctx)
	return nil
}

// 验证 request
func validateDeleteTenantRequest(req *pb.DeleteTenantRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
