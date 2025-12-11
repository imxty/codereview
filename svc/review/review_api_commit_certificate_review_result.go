package review

import (
	"context"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

var (
	notificationTitle = map[bool]string{
		true:  "商户审核资质通过",
		false: "商户审核资质未通过",
	}
)

const (
	ErrReviewTenantNotExist = 5303
)

// CommitCertificateReviewResult
// 每次审核都是新的工单，直接添加result即可
func (s *ReviewAPIHandler) CommitCertificateReviewResult(ctx context.Context, req *pb.CommitCertificateReviewResultRequest, rsp *pb.CommitCertificateReviewResultResponse) error {

	// 1.验证request
	err := validateCommitCertificateReviewResultRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	now := time.Now().UTC()
	// 创建工单结果
	// 验证工单
	issue, err := s.reviewStore.GetReviewIssueById(ctx, req.GetReviewIssueId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if issue == nil {
		return errors.Errorf(ErrReviewIssueNotFound, "issue not found[%s]", req.GetReviewIssueId())
	}
	if issue.GetTargetType() != domain.ReviewTypeCertificate {
		return errors.Errorf(codes.InvalidRequest, "issue type must be certificate[%s]", req.GetReviewIssueId())
	}
	// 查看工单的状态是否正确
	if issue.GetReviewStatus() != domain.ReviewStatusReviewing {
		return errors.Errorf(codes.InvalidRequest, "issue status should be reviewing [%s]", req.GetReviewIssueId())
	}

	// 修改工单状态
	if req.GetPass() {
		// 修改状态
		err = s.reviewStore.PassReviewIssue(ctx, req.GetReviewIssueId(), issue.GetRev())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	} else {
		// 修改状态
		err = s.reviewStore.ReviewIssueFailed(ctx, req.GetReviewIssueId(), issue.GetRev())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}
	// 如果通过则终止工单状态，否则就待用户输入
	_, err = s.userAPI.UpdateTenantCertificateByReviewResult(ctx, &userpb.UpdateTenantCertificateByReviewResultRequest{
		// 商户 ID
		TenantId: issue.GetSubmitterTenantID(),
		// 审核状态 成功 / 失败
		Status: req.GetPass(),
		// 失败原因，成功为空
		FailReason: req.GetComment(),
	})
	if err != nil {
		return err
	}

	// 创建工单结果
	issueResult := &domain.ReviewResult{
		ReviewResultID: xid.New().String(),
		ReviewIssueID:  issue.GetReviewIssueID(),
		ReviewTime:     now,
		ReviewerUserID: req.GetReviewerUserId(),
		Result:         req.GetPass(),
		Comment:        req.GetComment(),
	}
	err = s.reviewStore.CommitReviewResult(ctx, issueResult)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 创建审核通知
	notification := &domain.ReviewNotification{
		ReviewNotificationID: xid.New().String(),
		OrganizationID:       issue.GetSubmitterOrganizationID(),
		TenantID:             issue.GetSubmitterTenantID(),
		NotificationTitle:    notificationTitle[req.GetPass()],
		TargetType:           issue.GetTargetType(),
		Result:               req.GetPass(),
		Comment:              req.GetComment(),
	}
	err = s.reviewStore.CreateReviewNotification(ctx, notification)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	return nil
}

// 验证request
func validateCommitCertificateReviewResultRequest(req *pb.CommitCertificateReviewResultRequest) error {
	if req.GetReviewIssueId() == "" {
		return gerr.New("issue_id should not be empty")
	}
	if req.GetReviewerUserId() == "" {
		return gerr.New("review_user_id should not be empty")
	}
	if !req.GetPass() && req.GetComment() == "" {
		return gerr.New("comment should not be empty")
	}
	return nil
}
