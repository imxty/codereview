package user

import (
	"context"
	gerr "errors"
	"sort"
	"strings"

	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	reportv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 查询组织商户数据导出请求
func (u *UserAPIHandler) SearchOrganizationTenantDownload(ctx context.Context, req *pb.SearchOrganizationTenantDownloadRequest, rsp *pb.SearchOrganizationTenantDownloadResponse) error {
	err := validateSearchOrganizationTenantDownloadRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取组织信息
	organization, err := u.userStore.GetOrganizationByOrganizationId(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if organization == nil {
		return errors.Errorf(ErrOrganizationNotExist, "organization[%s] not found", req.GetOrganizationId())
	}

	// 获取组织下的所有商户信息
	temp, err := u.userStore.ListOrganizationAuthTenants(ctx, req.GetOrganizationId())
	if err != nil {
		return err
	}
	tids := make([]string, 0)
	for _, t := range temp {
		tids = append(tids, t.GetTenantID())
	}
	// 没有商户，返回空
	if len(tids) < 1 {
		return nil
	}
	tenants, err := u.userStore.ListTenantEntity(ctx, tids)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获取常客信息
	customerSearchResp, err := u.customerAPI.SearchOrganizationTenantCompare(ctx, &customerv1.SearchOrganizationTenantCompareRequest{
		OrganizationId: req.GetOrganizationId(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
		TenantName:     req.GetTenantName(),
		Pagination: &customerv1.Pagination{
			Offset: 0,
			Size:   organization.GetTenantLimit(),
		},
	})
	if err != nil {
		return err
	}

	// 获取测量信息
	reportSearchResp, err := u.reportAPI.ListOrganizationTenantReportDownload(ctx, &reportv1.ListOrganizationTenantReportDownloadRequest{
		OrganizationId: req.GetOrganizationId(),
		TenantName:     req.GetTenantName(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
	})
	if err != nil {
		return err
	}

	// 按商户名归类
	tenantsInfoMap := make(map[string]domain.TenantEntityIntf)
	for _, v := range tenants {
		if v.GetStoreName() == "" {
			continue
		}
		if req.GetTenantName() != "" {
			if !strings.Contains(v.GetStoreName(), req.GetTenantName()) {
				continue
			}
		}
		tenantsInfoMap[v.GetStoreName()] = v
	}
	customerInfoMap := make(map[string]*customerv1.TenantCustomerCompare)
	for _, v := range customerSearchResp.GetTenants() {
		customerInfoMap[v.GetTenantName()] = v
	}
	reportInfoMap := make(map[string]*reportv1.TenantReportCount)
	for _, v := range reportSearchResp.GetCounts() {
		reportInfoMap[v.GetTenantName()] = v
	}

	results := make([]*pb.OrganizationTenantDownload, 0, len(tenantsInfoMap))
	for storeName, tenant := range tenantsInfoMap {
		report := reportInfoMap[storeName]
		customer := customerInfoMap[storeName]

		measurementCount := int32(0)
		if report != nil {
			measurementCount = report.GetMonthlyCustomerMeasurementCount()
		}

		customerCount := int32(0)
		if customer != nil {
			customerCount = customer.GetMonthCount()
		}

		results = append(results, &pb.OrganizationTenantDownload{
			TenantName:       storeName,
			ContactName:      tenant.GetContactName(),
			ContactPhone:     tenant.GetContactPhone(),
			MeasurementCount: measurementCount,
			CustomerCount:    customerCount,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].MeasurementCount > results[j].MeasurementCount
	})

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
