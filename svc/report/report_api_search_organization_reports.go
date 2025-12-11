package report

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SearchOrganizationReports 查询报告测量总次数
func (s *ReportAPIHandler) SearchOrganizationReports(ctx context.Context, req *pb.SearchOrganizationReportsRequest, rsp *pb.SearchOrganizationReportsResponse) error {
	err := validateSearchOrganizationReportsRequest(req)
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
		// 如果没有tenantId直接返回即可
		if len(tids) == 0 {
			return nil
		}
	}

	// 查询报告
	rCounts, counts, err := s.reportStore.SearchOrganizationReports(ctx, req.GetOrganizationId(), tids, req.GetStartTime().AsTime().UTC(), req.GetEndTime().AsTime().UTC(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	dCounts := make([]*pb.ReportCount, len(rCounts))
	tid := make(map[string]bool)
	rtids := []string{}
	oid := make(map[string]bool)
	roids := []string{}
	for _, v := range rCounts {
		// 统计tenantID和organizationID
		if !oid[v.GetOrganizationID()] {
			oid[v.GetOrganizationID()] = true
			roids = append(roids, v.GetOrganizationID())
		}
		if !tid[v.GetTenantID()] {
			tid[v.GetTenantID()] = true
			rtids = append(rtids, v.GetTenantID())
		}
	}
	// 查询组织名称和商户名称
	getOrganizationsRsp, err := s.userAPI.BatchGetOrganizationNamesByIDs(ctx, &userv1.BatchGetOrganizationNamesByIDsRequest{
		OrganizationId: roids,
	})
	if err != nil {
		return err
	}
	getTenantsRsp, err := s.userAPI.BatchGetTenantNamesByIDs(ctx, &userv1.BatchGetTenantNamesByIDsRequest{
		TenantIds: rtids,
	})
	if err != nil {
		return err
	}
	tenantsName := getTenantsRsp.GetTenantNames()
	organizationsName := getOrganizationsRsp.GetOrganizations()
	// 返回数据
	for k, v := range rCounts {
		dCounts[k] = &pb.ReportCount{
			// 组织名称
			OrganizationName: organizationsName[v.GetOrganizationID()],
			// 商户名称
			TenantName: tenantsName[v.GetTenantID()],
			// 体验测量次数
			TempCustomerReportCount: int32(v.GetTempCount()),
			// vip测量次数
			CustomerReportCount: int32(v.GetCustomerCount()),
			// 开始时间
			StartTime: timestamppb.New(v.GetDate()),
		}
	}
	rsp.ReportCount = dCounts
	rsp.TotalCount = int32(counts)
	return nil
}

// 验证request
func validateSearchOrganizationReportsRequest(req *pb.SearchOrganizationReportsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	st := req.GetStartTime().AsTime()
	et := req.GetEndTime().AsTime()
	if st.After(et) {
		return gerr.New("invalid  time")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
