package h5

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/h5/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *H5APIHandler) GetCustomerTrendency(ctx context.Context, req *pb.GetCustomerTrendencyRequest, rsp *pb.GetCustomerTrendencyResponse) error {
	err := validateGetCustomerTrendencyRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取公开分享的慧脉宝药报告请求
	getRsp, err := s.reportAPI.GetCustomerTrendency(ctx, &reportpb.GetCustomerTrendencyRequest{
		// 租户id
		TenantId: req.GetTenantId(),
		// 常客ID
		CustomerId: req.GetCustomerId(),
		Days:       req.GetDays(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.CInfos = toAppCInfos(getRsp.GetCInfos())
	rsp.StressStateCount = toAppStressStateCount(getRsp.GetStressStateCount())
	return nil
}

// 验证request
func validateGetCustomerTrendencyRequest(req *pb.GetCustomerTrendencyRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomerId() == "" {
		return gerr.New("customer id should not be empty")
	}
	if req.GetDays() <= 0 {
		return gerr.New("invalid days")
	}
	return nil
}
