package review

import (
	"context"
	gerr "errors"
	"time"

	ptime "github.com/jinmukeji/huimaibao-service/pkg/time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

func (s *ReviewAPIHandler) CommitEntityCertificate(ctx context.Context, req *pb.CommitEntityCertificateRequest, rsp *pb.CommitEntityCertificateResponse) error {
	err := validateCommitEntityCertificateRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	now := time.Now().UTC()
	// 获取revision
	getRsp, err := s.userAPI.GetEntityRevision(ctx, &userv1.GetEntityRevisionRequest{
		RevisionId: req.GetTenantEntityRevisionId(),
	})
	if err != nil {
		return err
	}
	revision := getRsp.GetEntity()
	if revision == nil {
		return errors.Errorf(codes.InvalidRequest, "tenant entity revision not found [%s,%s]", req.GetTenantId(), req.GetTenantEntityRevisionId())
	}

	// 创建review issue
	issue := &domain.ReviewIssue{
		ReviewIssueID:           xid.New().String(),
		SubmitterOrganizationID: req.GetOrganizationId(),
		SubmitterTenantID:       req.GetTenantId(),
		SubmitTime:              now,
		AcceptanceTime:          ptime.InitTime,
		ReviewStatus:            domain.ReviewStatusPending,
		TargetType:              domain.ReviewTypeCertificate,
		TargetRevID:             req.GetTenantEntityRevisionId(),
	}
	err = s.reviewStore.CreateReviewIssue(ctx, issue)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证request
func validateCommitEntityCertificateRequest(req *pb.CommitEntityCertificateRequest) error {
	if req.GetIsOrganization() && req.GetOrganizationId() == "" {
		return gerr.New("organization should not be empty")
	}
	if req.GetTenantEntityRevisionId() == "" {
		return gerr.New("entity revision id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
