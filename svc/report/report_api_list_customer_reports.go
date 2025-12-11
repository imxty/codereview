package report

import (
	"context"
	"encoding/json"
	gerr "errors"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/huimaibao-service/svc/summary"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ReportMinCount 获取常客的体质报告所需报告数量
	ReportMinCount = 5
)

type PhysiqueJson struct {
	Physique string `json:"physique"`
	Score    int32  `json:"score"`
}

// ListCustomerReports 获取常客报告列表
func (s *ReportAPIHandler) ListCustomerReports(ctx context.Context, req *pb.ListCustomerReportsRequest, rsp *pb.ListCustomerReportsResponse) error {

	// 1.验证request
	err := validateListCustomerReportsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取常客信息
	customer, err := s.customerAPI.GetCustomer(ctx, &customerpb.GetCustomerRequest{
		TenantId:   req.GetTenantId(),
		CustomerId: req.GetCustomerId(),
	})
	if err != nil {
		return err
	}
	rsp.Customer = toProtoCustomer(customer.GetCustomer())
	// 获取常客最近测量的5笔报告
	reports, err := s.reportStore.ListCustomerLastReports(ctx, req.GetTenantId(), req.GetCustomerId(), ReportMinCount)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 报告数量不足返回错误
	if len(reports) >= ReportMinCount {
		// 返回体质map
		summary, err := getConstitution(reports)
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
		physique := []*PhysiqueJson{}
		for k, v := range summary {
			physique = append(physique, &PhysiqueJson{
				Physique: k,
				Score:    v,
			})
		}
		res, err := json.Marshal(physique)
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
		rsp.Physique = string(res)
	}

	// 获取报告
	if req.GetReportId() != "" {
		// 查询报告
		report, err := s.reportStore.GetCustomerReport(ctx, req.GetTenantId(), req.GetReportId(), req.GetCustomerId())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		if report == nil {
			return nil
		}
		// 转化格式
		simpleReports, err := toProtoSimpleReport(report)
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
		rsp.Reports = []*pb.SimpleHealthReport{simpleReports}
		rsp.TotalCount = 1
	} else {
		// 检测时间
		if req.GetTimeRange() == nil {
			// 如果没有时间就是查询全部
			reports, reportSize, err := s.reportStore.ListReportsWithoutTime(ctx, req.GetTenantId(), req.GetCustomerId(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
			// 转化数据格式
			simpleReports := make([]*pb.SimpleHealthReport, len(reports))
			for k, report := range reports {
				simpleReports[k], err = toProtoSimpleReport(report)
				if err != nil {
					return errors.Error(codes.InvalidOperation, err.Error())
				}
			}
			rsp.Reports = simpleReports
			rsp.TotalCount = int32(reportSize)
			return nil
		}
		timeRange := req.GetTimeRange()
		startTime := timeRange.GetStartTime().AsTime()
		endTime := timeRange.GetEndTime().AsTime()
		if startTime.After(endTime) {
			return errors.Error(codes.InvalidRequest, "end_time should not earlier than start_time")
		}
		// 查询报告
		reports, reportSize, err := s.reportStore.ListReports(ctx, req.GetTenantId(), req.GetCustomerId(), startTime, endTime, req.GetPagination().GetOffset(), req.GetPagination().GetSize())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 转化数据格式
		simpleReports := make([]*pb.SimpleHealthReport, len(reports))
		for k, report := range reports {
			simpleReports[k], err = toProtoSimpleReport(report)
			if err != nil {
				return errors.Errorf(codes.InvalidOperation, "failed to get simple report[%s]", err.Error())
			}
		}
		rsp.Reports = simpleReports
		rsp.TotalCount = int32(reportSize)
	}

	return nil
}

// 验证request
func validateListCustomerReportsRequest(req *pb.ListCustomerReportsRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomerId() == "" {
		return gerr.New("customer_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}

// getConstitution 获取体质占比
func getConstitution(reports []domain.ReportIntf) (map[string]int32, error) {
	res := make(map[string]int32)
	var rLen int32 = 0
	// 获取体质出现的次数
	for _, v := range reports {
		if v.GetPhysicalDialectics() == "" {
			continue
		}
		// 新报告
		var physique *domain.PhysiqueDialectics
		err := json.Unmarshal([]byte(v.GetPhysicalDialectics()), &physique)
		if err != nil {
			return nil, err
		}
		if len(physique.LookUps) == 0 {
			continue
		}
		lookup := physique.LookUps
		dialect := lookup[0]
		label := summary.PhysiqueMap[dialect.Key].Label
		res[label] = res[label] + 1
		rLen++
	}
	// 转化为百分比
	for k, v := range res {
		res[k] = v * 100 / rLen
	}
	return res, nil
}
