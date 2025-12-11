package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *BossAPIHandler) SearchStaffReportsCount(ctx context.Context, req *pb.SearchStaffReportsCountRequest, rsp *pb.SearchStaffReportsCountResponse) error {
	err := validateSearchStaffReportsCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 发送请求
	searchRsp, err := s.reportAPI.SearchStaffReportsCount(ctx, &reportpb.SearchStaffReportsCountRequest{
		// 商户ID
		TenantId: req.GetTenantId(),
		// 状态
		IsActivated: req.GetIsActivated(),
		// 开始时间
		StartTime: req.GetStartTime(),
		// 结束时间
		EndTime: req.GetEndTime(),
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.CustomerCount = searchRsp.GetCustomerCount()
	rsp.TempCustomerCount = searchRsp.GetTempCustomerCount()
	return nil
}

// 验证request
func validateSearchStaffReportsCountRequest(req *pb.SearchStaffReportsCountRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	st := req.GetStartTime().AsTime()
	et := req.GetEndTime().AsTime()
	if st.After(et) {
		return gerr.New("invalid  time")
	}
	return nil
}
