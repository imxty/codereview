package user

import (
	"context"
	gerr "errors"

	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) GetEntityOverdueSummary(ctx context.Context, req *pb.GetEntityOverdueSummaryRequest, rsp *pb.GetEntityOverdueSummaryResponse) error {
	err := validateGetEntityOverdueSummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取待复查人数
	getRsp, err := u.customerAPI.GetOverdueCount(ctx, &customerv1.GetOverdueCountRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return err
	}

	rsp.OverdueCustomerTotalCount = getRsp.GetTotalCount()

	return nil
}

// 验证 request
func validateGetEntityOverdueSummaryRequest(req *pb.GetEntityOverdueSummaryRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
