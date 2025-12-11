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

func (s *OrganizationAPIHandler) SubmitTreatmentToReview(ctx context.Context, req *pb.SubmitTreatmentToReviewRequest, rsp *pb.SubmitTreatmentToReviewResponse) error {
	err := validateSubmitTreatmentToReviewRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.productAPI.SubmitTreatmentToReview(ctx, &productv1.SubmitTreatmentToReviewRequest{
		OrganizationId: req.GetOrganizationId(),
		TreatmentId:    req.GetTreatmentId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

// 验证request
func validateSubmitTreatmentToReviewRequest(req *pb.SubmitTreatmentToReviewRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment id should not be empty")
	}
	return nil
}
