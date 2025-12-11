package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 查询商户名是否存在
func (u *UserAPIHandler) CheckTenantTreatmentUsing(ctx context.Context, req *pb.CheckTenantTreatmentUsingRequest, rsp *pb.CheckTenantTreatmentUsingResponse) error {
	err := validateCheckTeantTreatmentUsingRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	d, err := u.userStore.GetTreatmentTenant(ctx, req.GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.Using = d != nil

	return nil
}

func validateCheckTeantTreatmentUsingRequest(req *pb.CheckTenantTreatmentUsingRequest) error {
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment_id should not be empty")
	}
	return nil
}
