package review

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ReviewAPIHandler) ListReviewResultByTargetIds(ctx context.Context, req *pb.ListReviewResultByTargetIdsRequest, rsp *pb.ListReviewResultByTargetIdsResponse) error {
	err := validateListReviewResultByTargetIdsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取targetIssue
	issues, err := s.reviewStore.ListReviewIssueByTargetId(ctx, req.GetTargetId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	issueIds := make([]string, len(issues))
	for k, v := range issues {
		issueIds[k] = v.GetReviewIssueID()
	}
	issuesMap := make(map[string]domain.ReviewIssueIntf)
	for _, v := range issues {
		issuesMap[v.GetTargetRevID()] = v
	}
	// 查询result
	results, err := s.reviewStore.ListReviewResultByIssueId(ctx, issueIds)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	resultsMap := make(map[string]domain.ReviewResultIntf)
	for _, v := range results {
		resultsMap[v.GetReviewIssueID()] = v
	}
	reviewResult := make(map[string]*pb.ReviewResult)
	// 返回结果
	for _, v := range req.GetTargetId() {
		// 获取issue
		if issue, ok := issuesMap[v]; ok {
			if result, exist := resultsMap[issue.GetReviewIssueID()]; exist {
				// 只返回有result的数据
				reviewResult[v] = &pb.ReviewResult{
					// 审核类型
					ReviewType: toProtoReviewType(issue.GetTargetType()),
					// 审核状态
					ReiewStatus: toProtoReviewStatus(issue.GetReviewStatus()),
					// 如果审核失败,失败原因
					FailReason: result.GetComment(),
				}
			}
		}
	}

	rsp.Results = reviewResult
	return nil
}

func validateListReviewResultByTargetIdsRequest(req *pb.ListReviewResultByTargetIdsRequest) error {
	if len(req.GetTargetId()) == 0 {
		return gerr.New("target ids should not be empty")
	}
	return nil
}
