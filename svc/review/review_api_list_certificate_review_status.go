package review

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ReviewAPIHandler) ListCertificateReviewStatus(ctx context.Context, req *pb.ListCertificateReviewStatusRequest, rsp *pb.ListCertificateReviewStatusResponse) error {
	err := validateListCertificateReviewStatusRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 查询tenant的最新资质issue
	tenantIssues := make(map[string]domain.ReviewIssueIntf)
	issues, err := s.reviewStore.ListTenantCertificateIssues(ctx, req.GetTenantIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	issuesId := make([]string, len(issues))
	for k, v := range issues {
		tenantIssues[v.GetSubmitterTenantID()] = v
		issuesId[k] = v.GetReviewIssueID()
	}
	// 如果有商户没有工单说明传入的id有问题
	if len(tenantIssues) != len(req.GetTenantIds()) {
		return errors.Error(codes.InvalidRequest, "tenant issue not found")
	}
	// 查询issue相关结果
	result, err := s.reviewStore.ListReviewResultByIssueId(ctx, issuesId)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	resultMap := make(map[string]domain.ReviewResultIntf)
	for _, v := range result {
		resultMap[v.GetReviewIssueID()] = v
	}
	reviewResult := make(map[string]*pb.ReviewResult)
	for _, v := range req.GetTenantIds() {
		issue := tenantIssues[v]
		if rr, ok := resultMap[issue.GetReviewIssueID()]; !ok {
			// 如果还没有审核结果
			reviewResult[v] = &pb.ReviewResult{
				// 审核类型
				ReviewType: toProtoReviewType(issue.GetTargetType()),
				// 审核状态
				ReiewStatus: toProtoReviewStatus(issue.GetReviewStatus()),
			}
		} else {
			// 如果还没有审核结果
			reviewResult[v] = &pb.ReviewResult{
				// 审核类型
				ReviewType: toProtoReviewType(issue.GetTargetType()),
				// 审核状态
				ReiewStatus: toProtoReviewStatus(issue.GetReviewStatus()),
				// 如果审核失败,失败原因
				FailReason: rr.GetComment(),
			}
		}
	}

	rsp.Results = reviewResult

	return nil
}

// 验证request
func validateListCertificateReviewStatusRequest(req *pb.ListCertificateReviewStatusRequest) error {
	if len(req.GetTenantIds()) == 0 {
		return gerr.New("tenant ids should not be empty")
	}
	return nil
}
