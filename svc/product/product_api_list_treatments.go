package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ProductAPIHandler) ListTreatments(ctx context.Context, req *pb.ListTreatmentsRequest, rsp *pb.ListTreatmentsResponse) error {
	// 验证request
	err := validateListTreatmentsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 通过租户id查询方案版本信息
	var dTreatmentRevs []domain.TreatmentRevIntf
	// 查询最新的方案
	dTreatmentRevs, err = s.productStore.ListLatestTreatmentRevsWithStatus(ctx, req.GetOrganizationId(), req.GetIsPublish())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 如果没有方案直接返回即可
	if len(dTreatmentRevs) == 0 {
		return nil
	}
	targetIds := make([]string, len(dTreatmentRevs))
	for k, v := range dTreatmentRevs {
		targetIds[k] = v.GetTreatmentRevID()
	}
	// 获取审核结果
	getRsp, err := s.reviewAPI.ListReviewResultByTargetIds(ctx, &reviewpb.ListReviewResultByTargetIdsRequest{
		TargetId: targetIds,
	})
	if err != nil {
		return err
	}
	reviewResults := getRsp.GetResults()

	pbTreatments := make([]*pb.Treatment, len(dTreatmentRevs))

	// 转换
	for i, v := range dTreatmentRevs {
		if result, ok := reviewResults[v.GetTreatmentRevID()]; !ok {
			// 没有结果(审核不是失败的状态)
			if v.GetTreatmentStatus() == domain.TreatmentStatusApproved {
				pbTreatments[i] = transferToProtoTreatmentInfo(v, true, "")
				continue
			}
			pbTreatments[i] = transferToProtoTreatmentInfo(v, false, "")
		} else {
			// 没有结果(审核不是失败的状态)
			if v.GetTreatmentStatus() == domain.TreatmentStatusApproved {
				pbTreatments[i] = transferToProtoTreatmentInfo(v, true, "")
				continue
			}
			pbTreatments[i] = transferToProtoTreatmentInfo(v, false, result.GetFailReason())
		}
	}
	rsp.Treatments = pbTreatments
	return nil
}

// 验证request
func validateListTreatmentsRequest(req *pb.ListTreatmentsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	return nil
}
