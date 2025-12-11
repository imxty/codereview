package review

import (
	"context"
	gerr "errors"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// ErrReviewIssueNotFound
	ErrReviewIssueNotFound = 5026
)

// AcceptReviewIssue 接受工单
// 1.查询工单，判断状态
// 2.让工单变成该员工并修改工单状态
// 3.返回工单内容
func (s *ReviewAPIHandler) AcceptReviewIssue(ctx context.Context, req *pb.AcceptReviewIssueRequest, rsp *pb.AcceptReviewIssueResponse) error {

	// 1.验证request
	err := validateAcceptReviewIssueRequest(req)
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

	// 查看工单状态是否可以被接单
	if issue.GetReviewStatus() != domain.ReviewStatusPending {
		return errors.Errorf(codes.InvalidOperation, "fail to accept review issue[%s]", req.GetReviewIssueId())
	}

	// 接单
	err = s.reviewStore.AcceptReviewIssue(ctx, issue.GetReviewIssueID(), req.GetReviewerUserId(), issue.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.ReviewIssueId = req.GetReviewIssueId()

	// 查询当前工单的内容
	// 如果是接受的证书审核
	if issue.GetTargetType() == domain.ReviewTypeCertificate {
		// 查询entity_revision
		revisionRsp, err := s.userAPI.GetEntityRevision(ctx, &userv1.GetEntityRevisionRequest{
			RevisionId: issue.GetTargetRevID(),
		})
		if err != nil {
			return err
		}
		revision := revisionRsp.GetEntity()
		entity := &pb.TenantEntity{
			// 租户ID
			TenantId: revision.GetTenantId(),
			// 商铺名称
			EntityName: revision.GetName(),
			// 地址
			Address: &pb.Address{
				// 省
				Province: revision.GetAddress().GetProvince(),
				// 市
				City: revision.GetAddress().GetCity(),
				// 区
				District: revision.GetAddress().GetDistrict(),
				// 街道
				Street: revision.GetAddress().GetStreet(),
			},
			// 联系人姓名
			ContactName: revision.GetContactName(),
			// 联系人电话
			ContactPhone: revision.GetContactPhone(),
			// 营业执照
			BusinessLicenseUrl: revision.GetBusinessLicenseUrl(),
			// 社会信用代码
			SocialCreditCode: revision.GetSocialCreditCode(),
			LogoUrl:          revision.GetLogoUrl(),
		}
		// 返回数据
		rsp.ReviewType = pb.ReviewType_REVIEW_TYPE_CERTIFICATE
		rsp.Entity = entity
		rsp.SubmitTime = timestamppb.New(issue.GetSubmitTime())
		rsp.SubmitterTenantId = issue.GetSubmitterTenantID()
		rsp.SubmitterOrganizationId = issue.GetSubmitterOrganizationID()
	} else {
		// 获取treatment的信息
		getTreatment, err := s.productAPI.GetTreatmentRevByTreatmentRevId(ctx, &productpb.GetTreatmentRevByTreatmentRevIdRequest{
			OrganizationId: issue.GetSubmitterOrganizationID(),
			TreatmentRevId: issue.GetTargetRevID(),
		})
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 转换 treatment -> SymptomProductMap
		productMap, err := transferTreatmentToSymptomProductMap(getTreatment.GetTreatment())
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
		rsp.ProductMap = productMap
		rsp.SubmitTime = timestamppb.New(issue.GetSubmitTime())
		rsp.SubmitterTenantId = issue.GetSubmitterTenantID()
		rsp.SubmitterOrganizationId = issue.GetSubmitterOrganizationID()
		return nil
	}
	return nil
}

// 验证request
func validateAcceptReviewIssueRequest(req *pb.AcceptReviewIssueRequest) error {
	if req.GetReviewIssueId() == "" {
		return gerr.New("review_issued_id should not be empty")
	}
	if req.GetReviewerUserId() == "" {
		return gerr.New("review_user_id should not be empty")
	}
	return nil
}

// transferTreatmentToSymptomProductMap
func transferTreatmentToSymptomProductMap(tr *productpb.Treatment) ([]*pb.SymptomProductMap, error) {
	// 分别获取三种方案配置类型列表
	// 风险疾病方案配置
	rd := tr.GetTreatmentItemsRiskyDisease()
	// 脏腑辩证方案配置
	dd := tr.GetTreatmentItemsDirtyDialectic()
	// 理疗方案配置
	pt := tr.GetTreatmentItemsPhysicalTherapy()
	// 体质方案配置
	ptd := tr.GetTreatmentItemsPhysicalDialectics()
	// 合成一个列表
	list := append(append(append(rd, dd...), pt...), ptd...)
	res := make([]*pb.SymptomProductMap, len(list))
	for i, v := range list {
		spm, err := transferTreatmentItemToSymptomProductMap(v)

		if err != nil {
			return nil, err
		}
		res[i] = spm
	}
	return res, nil
}

// TreatmentItem -> SymptomProductMap
func transferTreatmentItemToSymptomProductMap(tr *productpb.TreatmentItem) (*pb.SymptomProductMap, error) {
	products := make([]*pb.Product, len(tr.GetProducts()))
	for i, v := range tr.GetProducts() {
		p, err := toReviewProtoProductFromProduct(v)
		if err != nil {
			return nil, err
		}
		products[i] = p
	}
	return &pb.SymptomProductMap{
		Symptom:  tr.GetSymptom(),
		Products: products,
	}, nil
}
