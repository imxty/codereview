package report

import (
	"context"
	gerr "errors"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// SearchReports 查询报告
func (s *ReportAPIHandler) SearchReports(ctx context.Context, req *pb.SearchReportsRequest, rsp *pb.SearchReportsResponse) error {
	err := validateSearchReportsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 如果有常客数据则查询
	customerId := ""
	if req.GetCustomerName() != "" || req.GetCustomerPhone() != "" {
		searchRsp, err := s.customerAPI.SearchCustomerIDByCustomerInfo(ctx, &customerpb.SearchCustomerIDByCustomerInfoRequest{
			// 客户名称
			Name: req.GetCustomerName(),
			// 手机号
			Phone: req.GetCustomerPhone(),
		})
		if err != nil {
			return err
		}
		customerId = searchRsp.GetCustomerId()
	}
	// 如果有商户名称
	var tids []string
	if req.GetTenantName() != "" {
		searchRsp, err := s.userAPI.SearchTenantsByName(ctx, &userpb.SearchTenantsByNameRequest{
			TenantName: req.GetTenantName(),
		})
		if err != nil {
			return err
		}
		tids = searchRsp.GetTenantIds()
	}
	if req.GetTenantId() != "" {
		tids = append(tids, req.GetTenantId())
	}
	// 构建查询参数
	customerType := domain.CustomerTypeBoth
	switch req.GetCustomerType() {
	case pb.CustomerType_CUSTOMER_TYPE_CUSTOMER:
		customerType = domain.CustomerTypeCustomer
	case pb.CustomerType_CUSTOMER_TYPE_TEMP:
		customerType = domain.CustomerTypeTemp
	}

	// 查询报告
	// 获取要查询的key flag
	ddflag, err := getDDKeyFlag(req.GetDirtyDialectics())
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	reports, count, err := s.reportStore.SearchReports(ctx, req.GetReportId(), customerType, customerId, tids, ddflag, req.GetOrganizationId(),
		req.GetStartTime().AsTime(), req.GetEndTime().AsTime(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回数据
	rsp.TotalCount = count
	reportSummary := make([]*pb.SummaryReport, len(reports))
	// 查询ID信息
	tid := make(map[string]bool)
	rtids := []string{}
	cid := make(map[string]bool)
	rcids := []string{}
	for _, v := range reports {
		if v.GetIsCustomer() && !cid[v.GetCustomerID()] {
			cid[v.GetCustomerID()] = true
			rcids = append(rcids, v.GetCustomerID())
		}
		if !tid[v.GetTenantID()] {
			tid[v.GetTenantID()] = true
			rtids = append(rtids, v.GetTenantID())
		}
	}

	// 商户 ID 对应的商户名
	tenantNames := make(map[string]string)
	customerInfo := make(map[string]*customerpb.Customer)
	// 获取常客信息和商户信息
	if len(rtids) != 0 {
		getRsp, err := s.userAPI.BatchGetTenantNamesByIDs(ctx, &userpb.BatchGetTenantNamesByIDsRequest{
			TenantIds: rtids,
		})
		if err != nil {
			return err
		}
		tenantNames = getRsp.GetTenantNames()
	}
	if len(rcids) != 0 {
		getRsp, err := s.customerAPI.BatchGetCustomersByIds(ctx, &customerpb.BatchGetCustomersByIdsRequest{
			CustomerIds:    rcids,
			IncludeDeleted: true,
		})
		if err != nil {
			return err
		}
		customerInfo = getRsp.GetCustomers()
	}

	// 数据转化
	for k, v := range reports {
		reportSummary[k], err = toReportSummary(v)
		if err != nil {
			return errors.Errorf(codes.InvalidOperation, "failed to get report summary[%s]", err.Error())
		}
		if v.GetIsCustomer() {
			reportSummary[k].CustomerName = customerInfo[v.GetCustomerID()].GetNickname()
			reportSummary[k].CustomerPhone = customerInfo[v.GetCustomerID()].GetPhone()
			reportSummary[k].IsCustomer = true
		}
		reportSummary[k].TenantName = tenantNames[v.GetTenantID()]
	}

	rsp.Reports = reportSummary
	return nil
}

// 验证request
func validateSearchReportsRequest(req *pb.SearchReportsRequest) error {
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	if req.GetStartTime().AsTime().After(req.GetEndTime().AsTime()) {
		return gerr.New("invalid time")
	}
	// 不能同时存在
	if req.GetTenantName() != "" && req.GetTenantId() != "" {
		return gerr.New("invalid tenant search condition")
	}
	return nil
}
