package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 修改方案状态请求
func (s *ProductAPIHandler) ChangeTreatmentRevStatus(ctx context.Context, req *pb.ChangeTreatmentRevStatusRequest, rsp *pb.ChangeTreatmentRevStatusResponse) error {
	err := validateChangeTreatmentRevStatusRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取方案
	tr, err := s.productStore.GetTreatmentRevByTreatmentRevId(ctx, req.GetTreatmentRevId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tr == nil {
		return errors.Errorf(ErrGetTreatmentFailed, "treatment rev[%s] not found", req.GetTreatmentRevId())
	}
	status, err := toDomainTreatmentStatus(req.GetStatus())
	if err != nil {
		return errors.Errorf(codes.InvalidRequest, "invalid status[%s]", err.Error())
	}
	// 修改treatment rev状态
	err = s.productStore.ChangeTreatmentRevStatus(ctx, req.GetTreatmentRevId(), tr.GetRev(), status)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	return nil
}

func validateChangeTreatmentRevStatusRequest(req *pb.ChangeTreatmentRevStatusRequest) error {
	if req.GetTreatmentRevId() == "" {
		return gerr.New("treatment rev id should not be empty")
	}
	if req.GetStatus() == pb.TreatmentStatus_TREATMENT_STATUS_INVALID || req.GetStatus() == pb.TreatmentStatus_TREATMENT_STATUS_UNSET {
		return gerr.New("invalid status " + req.GetStatus().String())
	}
	return nil
}
