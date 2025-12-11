package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ProductAPIHandler) UpdateTreatmentStatus(ctx context.Context, req *pb.UpdateTreatmentStatusRequest, rsp *pb.UpdateTreatmentStatusResponse) error {
	// 验证request
	err := validateUpdateTreatmentStatusRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取方案
	treatmentRev, err := s.productStore.GetTreatmentRevByTreatmentRevId(ctx, req.GetTreatmentRevId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if treatmentRev == nil {
		return errors.Errorf(ErrGetTreatmentFailed, "get treatment rev failed[%s]", req.GetTreatmentRevId())
	}
	// 获取需要改变的方案状态
	status, err := toDomainTreatmentStatus(req.GetTreatmentStatus())
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	err = s.productStore.ChangeTreatmentRevStatus(ctx, treatmentRev.GetTreatmentRevID(), treatmentRev.GetRev(), status)
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "change treatment status failed [%s]", err.Error())
	}
	return nil
}

// 验证request
func validateUpdateTreatmentStatusRequest(req *pb.UpdateTreatmentStatusRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentRevId() == "" {
		return gerr.New("treatment rev id should not be empty")
	}
	if _, err := toDomainTreatmentStatus(req.GetTreatmentStatus()); err != nil {
		return err
	}
	return nil
}
