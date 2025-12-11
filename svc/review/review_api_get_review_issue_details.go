package review

import (
	"context"
	gerr "errors"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ReviewAPIHandler) GetReviewIssueDetails(ctx context.Context, req *pb.GetReviewIssueDetailsRequest, rsp *pb.GetReviewIssueDetailsResponse) error {
	// 1.验证request
	err := validateGetReviewIssueDetails(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询工单
	issue, err := s.reviewStore.GetReviewIssueById(ctx, req.GetReviewIssueId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if issue == nil {
		return errors.Error(ErrReviewIssueNotFound, "review issue not found")
	}
	if issue.GetReviewerUserID() != req.GetReviewerUserId() {
		return errors.Errorf(codes.InvalidOperation, "review issue[%s] not belong to this user", req.GetReviewIssueId())
	}

	// 返回工单内容
	// 如果是接受的证书审核
	if issue.GetTargetType() == domain.ReviewTypeCertificate {
		// 查询entity_revision
		getRsp, err := s.userAPI.GetEntityRevision(ctx, &userpb.GetEntityRevisionRequest{
			RevisionId: issue.GetTargetRevID(),
		})
		if err != nil {
			return err
		}
		// 返回数据
		rsp.ReviewType = pb.ReviewType_REVIEW_TYPE_CERTIFICATE
		rsp.Entity = toProtoTenantEntity(getRsp.GetEntity())
	} else {
		// 获取treatment的信息
		getTreatmentRev, err := s.productAPI.GetTreatmentRevByTreatmentRevId(ctx, &productpb.GetTreatmentRevByTreatmentRevIdRequest{
			OrganizationId: issue.GetSubmitterOrganizationID(),
			TreatmentRevId: issue.GetTargetRevID(),
		})
		if err != nil {
			return err
		}
		// 转换 treatment -> SymptomProductMap
		productMap, err := transferTreatmentToSymptomProductMap(getTreatmentRev.GetTreatment())
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}

		rsp.ReviewType = pb.ReviewType_REVIEW_TYPE_PRODUCT
		rsp.ProductMap = productMap
		return nil
	}
	return nil
}

// 验证request
func validateGetReviewIssueDetails(req *pb.GetReviewIssueDetailsRequest) error {
	if req.GetReviewIssueId() == "" {
		return gerr.New("review issue id should not be empty")
	}
	return nil
}
