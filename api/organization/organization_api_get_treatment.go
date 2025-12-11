package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	productv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *OrganizationAPIHandler) GetTreatment(ctx context.Context, req *pb.GetTreatmentRequest, rsp *pb.GetTreatmentResponse) error {
	err := validateGetTreatmentRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.productAPI.GetTreatment(ctx, &productv1.GetTreatmentRequest{
		OrganizationId: req.GetOrganizationId(),
		TreatmentId:    req.GetTreatmentId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.Treatment = toAppTreatment(getRsp.GetTreatment(), s.s3Domain)
	return nil
}

// 验证request
func validateGetTreatmentRequest(req *pb.GetTreatmentRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("report id should not be empty")
	}
	return nil
}
