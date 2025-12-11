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

func (s *BossAPIHandler) ListReviewIssues(ctx context.Context, req *pb.ListReviewIssuesRequest, rsp *pb.ListReviewIssuesResponse) error {
	err := validateListReviewIssuesRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	listRsp, err := s.reviewAPI.ListReviewIssues(ctx, &reviewv1.ListReviewIssuesRequest{
		// 搜索条件
		Status: toSvcReviewStatus(req.GetStatus()),
		// 是否搜索全部
		SearchAll: req.GetSearchAll(),
		// 分页
		Pagination: toReviewPagination(req.GetPagination()),
		// 审核人ID，用户ID
		ReviewerUserId:   req.GetReviewerUserId(),
		OrganizationName: req.GetOrganizationName(),
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	issues := listRsp.GetReviewIssues()
	pis := make([]*pb.ReviewIssue, len(issues))
	for k, v := range issues {
		pis[k] = toAppReviewIssue(v)
	}

	rsp.ReviewIssues = pis
	rsp.TotalCount = listRsp.GetTotalCount()
	return nil
}

// 验证request
func validateListReviewIssuesRequest(req *pb.ListReviewIssuesRequest) error {
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
