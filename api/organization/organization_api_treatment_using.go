package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 方案有没有被使用
func (s *OrganizationAPIHandler) TreatmentUsing(ctx context.Context, req *pb.TreatmentUsingRequest, rsp *pb.TreatmentUsingResponse) error {
	err := validateTreatmentUsingRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	return nil
}

// 验证request
func validateTreatmentUsingRequest(req *pb.TreatmentUsingRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment_id should not be empty")
	}
	return nil
}
