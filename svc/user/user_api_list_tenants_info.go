package user

import (
	"context"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取商户信息
func (s *UserAPIHandler) ListTenantsInfo(ctx context.Context, req *pb.ListTenantsInfoRequest, rsp *pb.ListTenantsInfoResponse) error {
	err := validateListTenantsInfoRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商户信息
	entities, err := s.userStore.ListTenantEntity(ctx, req.GetTenantIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 查询组织商户信息
	ots, err := s.userStore.ListTenants(ctx, req.GetTenantIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获取商户订阅信息
	ts, err := s.userStore.ListSubscriptionTimeline(ctx, req.GetTenantIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	ts_map := make(map[string]domain.SubscriptionTimelineIntf)
	for _, each := range ts {
		ts_map[each.GetTenantID()] = each
	}

	tenantStatus := make(map[string]pb.TenantStatus)
	tenantOrganizationId := make(map[string]string)
	for _, v := range ots {
		// 判断商户状态 未认证 / 待续期 / 使用中
		var tenant_status pb.TenantStatus
		if !v.GetIsActivated() {
			// 未认证
			tenant_status = pb.TenantStatus_TENANT_STATUS_UNAUTH
		} else {
			if t, ok := ts_map[v.GetTenantID()]; !ok {
				// 如果没有订阅信息，则是待续期状态
				tenant_status = pb.TenantStatus_TENANT_STATUS_PENDING
			} else {
				// 有订阅，判断订阅是否过期
				tenant_status = checkTenantSubscriptionTimeline(t)
			}
		}

		tenantStatus[v.GetTenantID()] = tenant_status
		tenantOrganizationId[v.GetTenantID()] = v.GetOrganizationID()
	}

	// 返回数据
	tenantMap := make(map[string]*pb.Entity)
	for _, v := range entities {
		tenantMap[v.GetTenantID()] = toProtoEntity(v)
		tenantMap[v.GetTenantID()].OrganizationId = tenantOrganizationId[v.GetTenantID()]
		tenantMap[v.GetTenantID()].TenantStatus = tenantStatus[v.GetTenantID()]
	}

	rsp.Tenants = tenantMap

	return nil
}

// 验证request
func validateListTenantsInfoRequest(req *pb.ListTenantsInfoRequest) error {
	if len(req.GetTenantIds()) == 0 {
		return gerr.New("tenants id should not be empty")
	}
	return nil
}

// 检查订阅有没有过期
func checkTenantSubscriptionTimeline(s domain.SubscriptionTimelineIntf) pb.TenantStatus {
	if time.Now().Before(s.GetEndTime().UTC()) {
		return pb.TenantStatus_TENANT_STATUS_USING
	}
	return pb.TenantStatus_TENANT_STATUS_PENDING
}
