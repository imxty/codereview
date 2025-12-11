package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 查询组织商户数据对比
func (s *OrganizationAPIHandler) SearchOrganizationTenantDownload(ctx context.Context, req *pb.SearchOrganizationTenantDownloadRequest, rsp *pb.SearchOrganizationTenantDownloadResponse) error {
	err := validateSearchOrganizationTenantDownloadRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchResp, err := s.userAPI.SearchOrganizationTenantDownload(ctx, &userv1.SearchOrganizationTenantDownloadRequest{
		OrganizationId: req.GetOrganizationId(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
		TenantName:     req.GetTenantName(),
	})
	if err != nil {
		return err
	}

	results := make([]*pb.OrganizationTenantDownload, len(searchResp.GetTenants()))
	for k, v := range searchResp.GetTenants() {
		results[k] = &pb.OrganizationTenantDownload{
			TenantName:       v.GetTenantName(),
			ContactName:      v.GetContactName(),
			ContactPhone:     v.GetContactPhone(),
			MeasurementCount: v.GetMeasurementCount(),
			CustomerCount:    v.GetCustomerCount(),
		}
	}

	rsp.Tenants = results

	return nil
}

func validateSearchOrganizationTenantDownloadRequest(req *pb.SearchOrganizationTenantDownloadRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be empty")
	}
	if req.GetEndTime() == nil {
		return gerr.New("end_time should not be empty")
	}
	return nil
}
