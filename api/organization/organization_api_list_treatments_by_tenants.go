package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	productv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取商户已分配方案列表
func (s *OrganizationAPIHandler) ListTreatmentsByTenants(ctx context.Context, req *pb.ListTreatmentsByTenantsRequest, rsp *pb.ListTreatmentsByTenantsResponse) error {
	err := validateListTreatmentsByTenantsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	if len(req.GetTenantIds()) < 1 {
		rsp.TenantTreatments = make(map[string]*pb.TenantTreatment)
		return nil
	}

	listRsp, err := s.userAPI.ListTreatmentsByTenantIDs(ctx, &userv1.ListTreatmentsByTenantIDsRequest{
		TenantIds: req.GetTenantIds(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 获取商户 ID 和方案 ID
	tenantIds := make([]string, 0, len(listRsp.GetTreatments()))
	treatmentIds := make([]string, 0, len(listRsp.GetTreatments()))
	for k, v := range listRsp.GetTreatments() {
		tenantIds = append(tenantIds, k)
		treatmentIds = append(treatmentIds, v)
	}

	// 获取方案名
	treatmentRsp, err := s.productAPI.ListTreatmentNameByIds(ctx, &productv1.ListTreatmentNameByIdsRequest{
		TreatmentIds: treatmentIds,
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 获取商户名
	tenantRsp, err := s.userAPI.BatchGetTenantNamesByIDs(ctx, &userv1.BatchGetTenantNamesByIDsRequest{
		TenantIds: tenantIds,
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 整合数据
	results := make(map[string]*pb.TenantTreatment)
	for k, v := range tenantRsp.GetTenantNames() {
		tId := listRsp.GetTreatments()[k]
		results[k] = &pb.TenantTreatment{
			TenantId:      k,
			TenantName:    v,
			TreatmentId:   tId,
			TreatmentName: treatmentRsp.GetTreatmentNames()[tId],
		}
	}

	rsp.TenantTreatments = results

	return nil
}

// 验证request
func validateListTreatmentsByTenantsRequest(req *pb.ListTreatmentsByTenantsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetTenantIds() == nil {
		return gerr.New("tenant_ids should not be empty")
	}
	return nil
}
