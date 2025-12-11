package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	reviewv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *BossAPIHandler) GetReviewDetails(ctx context.Context, req *pb.GetReviewDetailsRequest, rsp *pb.GetReviewDetailsResponse) error {
	err := validateGetReviewDetailsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 发起获取工单详情请求
	getDetailRsp, err := s.reviewAPI.GetReviewDetails(ctx, &reviewv1.GetReviewDetailsRequest{
		ReviewIssueId: req.GetReviewIssueId(),
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回结果
	rsp.Entity = toAppTenantEntity(s.s3Domain, getDetailRsp.GetEntity())
	appProductMap := make([]*pb.SymptomProductMap, len(getDetailRsp.GetProductMap()))
	for k, v := range getDetailRsp.GetProductMap() {
		appProductMap[k] = toAppProductMap(s.s3Domain, v)
	}
	rsp.ProductMap = appProductMap
	rsp.ReviewType = toApiReviewType(getDetailRsp.GetReviewType())
	rsp.ReviewIssueId = req.GetReviewIssueId()
	// 租户ID
	rsp.SubmitterTenantId = getDetailRsp.GetSubmitterTenantId()
	// 提审时间
	rsp.SubmitTime = getDetailRsp.GetSubmitTime()
	// 组织ID
	rsp.SubmitterOrganizationId = getDetailRsp.GetSubmitterOrganizationId()
	rsp.TreatmentName = getDetailRsp.GetTreatmentName()

	return nil
}

// 验证request
func validateGetReviewDetailsRequest(req *pb.GetReviewDetailsRequest) error {
	if req.GetReviewIssueId() == "" {
		return gerr.New("review issue id should not be empty")
	}
	return nil
}
