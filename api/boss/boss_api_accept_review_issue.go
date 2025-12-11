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

func (s *BossAPIHandler) AcceptReviewIssue(ctx context.Context, req *pb.AcceptReviewIssueRequest, rsp *pb.AcceptReviewIssueResponse) error {
	err := validateAcceptReviewIssueRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	acceptRsp, err := s.reviewAPI.AcceptReviewIssue(ctx, &reviewv1.AcceptReviewIssueRequest{
		ReviewIssueId:  req.GetReviewIssueId(),
		ReviewerUserId: req.GetReviewerUserId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	appProductMap := make([]*pb.SymptomProductMap, len(acceptRsp.GetProductMap()))
	for k, v := range acceptRsp.GetProductMap() {
		appProductMap[k] = toAppProductMap(s.s3Domain, v)
	}

	rsp.ReviewIssueId = acceptRsp.GetReviewIssueId()
	// 审核类型
	rsp.ReviewType = toApiReviewType(acceptRsp.GetReviewType())
	// 租户信息
	rsp.Entity = toAppTenantEntity(s.s3Domain, acceptRsp.GetEntity())
	// 症候商品
	rsp.ProductMap = appProductMap
	// 租户ID
	rsp.SubmitterTenantId = acceptRsp.GetSubmitterTenantId()
	// 提审时间
	rsp.SubmitTime = acceptRsp.GetSubmitTime()
	// 组织ID
	rsp.SubmitterOrganizationId = acceptRsp.GetSubmitterOrganizationId()
	return nil
}

// 验证request
func validateAcceptReviewIssueRequest(req *pb.AcceptReviewIssueRequest) error {
	if req.GetReviewIssueId() == "" {
		return gerr.New("review issue id should not be empty")
	}
	if req.GetReviewerUserId() == "" {
		return gerr.New("review user id should not be empty")
	}
	return nil
}
