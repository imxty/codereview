package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) ListTreatmentsByTenantIDs(ctx context.Context, req *pb.ListTreatmentsByTenantIDsRequest, rsp *pb.ListTreatmentsByTenantIDsResponse) error {
	err := validateListTreatmentsByTenantIDsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取方案
	listRsp, err := u.userStore.ListTreatmentsByTenantIDs(ctx, req.GetTenantIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 数据转换
	results := make(map[string]string)
	for _, each := range listRsp {
		results[each.GetTenantID()] = each.GetTreatmentID()
	}

	rsp.Treatments = results

	return nil
}

// 验证 request
func validateListTreatmentsByTenantIDsRequest(req *pb.ListTreatmentsByTenantIDsRequest) error {
	if req.GetTenantIds() == nil {
		return gerr.New("tenant_ids should not be empty")
	}
	if len(req.GetTenantIds()) < 1 {
		return gerr.New("tenant_ids should not be empty")
	}
	return nil
}
