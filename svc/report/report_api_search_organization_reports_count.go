package report

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// SearchOrganizationReportsCount 查询报告测量总次数
func (s *ReportAPIHandler) SearchOrganizationReportsCount(ctx context.Context, req *pb.SearchOrganizationReportsCountRequest, rsp *pb.SearchOrganizationReportsCountResponse) error {
	err := validateSearchOrganizationReportsCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 如果商户名称不为空则需要查询匹配的商户
	tids := []string{}
	if req.GetTenantName() != "" {
		getRsp, err := s.userAPI.SearchTenantByNameAndOrganizationID(ctx, &userv1.SearchTenantByNameAndOrganizationIDRequest{
			OrganizationId: req.GetOrganizationId(),
			TenantName:     req.GetTenantName(),
		})
		if err != nil {
			return nil
		}
		tids = getRsp.GetTenantIds()
	}

	// 查询报告
	cCount, tCount, err := s.reportStore.SearchOrganizationReportsCount(ctx, req.GetOrganizationId(), tids, req.GetStartTime().AsTime(), req.GetEndTime().AsTime())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.CustomerCount = int32(cCount)
	rsp.TempCustomerCount = int32(tCount)
	return nil
}

// 验证request
func validateSearchOrganizationReportsCountRequest(req *pb.SearchOrganizationReportsCountRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	st := req.GetStartTime().AsTime()
	et := req.GetEndTime().AsTime()
	if st.After(et) {
		return gerr.New("invalid  time")
	}
	return nil
}
