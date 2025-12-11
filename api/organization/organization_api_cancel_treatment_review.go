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

// CancelTreatmentReview
func (s *OrganizationAPIHandler) CancelTreatmentReview(ctx context.Context, req *pb.CancelTreatmentReviewRequest, rsp *pb.CancelTreatmentReviewResponse) error {
	err := validateCancelTreatmentReviewRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 发送请求
	_, err = s.productAPI.CancelTreatmentReview(ctx, &productv1.CancelTreatmentReviewRequest{
		OrganizationId: req.GetOrganizationId(),
		TreatmentId:    req.GetTreatmentId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateCancelTreatmentReviewRequest(req *pb.CancelTreatmentReviewRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment id should not be empty")
	}
	return nil
}
