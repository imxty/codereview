package review

import (
	"context"
	gerr "errors"
	"time"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

func (s *ReviewAPIHandler) CommitTreatmentReviewResult(ctx context.Context, req *pb.CommitTreatmentReviewResultRequest, rsp *pb.CommitTreatmentReviewResultResponse) error {
	// 验证request
	err := validateCommitTreatmentReviewResultRequest(req)
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
		return errors.Error(ErrReviewIssueNotFound, "issue not found")
	}
	// 查看工单是否被用户取消审核
	if issue.GetIsCancel() {
		rsp.CancelReview = true
		return nil
	}
	rsp.CancelReview = false

	// 查看工单的状态是否正确
	if issue.GetReviewStatus() != domain.ReviewStatusReviewing {
		return errors.Errorf(codes.InvalidRequest, "issue[%s] status should be reviewing", req.GetReviewIssueId())
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

	// 修改工单状态
	// 如果通过,修改方案状态，商品的状态为 待使用状态
	if req.GetPass() {
		// 修改工单状态为审核终止
		err = s.reviewStore.PassReviewIssue(ctx, req.GetReviewIssueId(), issue.GetRev())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}
	// 判断租户方案的状态
	// 审核通过则审核完成
	var treatmentStatus productpb.TreatmentStatus
	if req.GetPass() {
		treatmentStatus = productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED
	} else {
		treatmentStatus = productpb.TreatmentStatus_TREATMENT_STATUS_UNAPPROVED
	}

	// 更新租户方案审核状态
	_, err = s.productAPI.UpdateTreatmentStatus(ctx, &productpb.UpdateTreatmentStatusRequest{
		OrganizationId:  issue.GetSubmitterOrganizationID(),
		TreatmentRevId:  issue.GetTargetRevID(),
		TreatmentStatus: treatmentStatus,
	})
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 审核通过后 修改方案中所有商品的状态为 待使用状态
	if req.GetPass() {
		_, err = s.productAPI.ChangeProductStatusByTreatmentRevId(ctx, &productpb.ChangeProductStatusByTreatmentRevIdRequest{
			TreatmentRevId: issue.GetTargetRevID(),
			ProductStatus:  productpb.ProductStatus_PRODUCT_STATUS_TO_BE_USED,
		})
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	} else {
		// 没通过
		err = s.reviewStore.ReviewIssueFailed(ctx, req.GetReviewIssueId(), issue.GetRev())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 2.更新租户方案审核状态 未过审
		_, err = s.productAPI.UpdateTreatmentStatus(ctx, &productpb.UpdateTreatmentStatusRequest{
			OrganizationId:  issue.GetSubmitterOrganizationID(),
			TreatmentRevId:  issue.GetTargetRevID(),
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_UNAPPROVED,
		})
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}
	// 获取方案名称
	getTreatmentRevRsp, err := s.productAPI.GetTreatmentRevByTreatmentRevId(ctx, &productpb.GetTreatmentRevByTreatmentRevIdRequest{
		OrganizationId: issue.GetSubmitterOrganizationID(),
		TreatmentRevId: issue.GetTargetRevID(),
	})
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 创建审核通知
	notification := &domain.ReviewNotification{
		ReviewNotificationID: xid.New().String(),
		OrganizationID:       issue.GetSubmitterOrganizationID(),
		TenantID:             issue.GetSubmitterTenantID(),
		NotificationTitle:    getTreatmentRevRsp.GetTreatment().GetTreatmentName(),
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
func validateCommitTreatmentReviewResultRequest(req *pb.CommitTreatmentReviewResultRequest) error {
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
