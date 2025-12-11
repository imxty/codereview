package user

import (
	"context"
	gerr "errors"
	"slices"

	productv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrTenantNotFound
	ErrTenantNotFound = 5004
)

func (u *UserAPIHandler) SearchTenants(ctx context.Context, req *pb.SearchTenantsRequest, rsp *pb.SearchTenantsResponse) error {
	err := validateSearchTenantsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 如果组织为空，返回空数据
	if req.GetOrganizationId() == "" {
		rsp.Tenants = []*pb.TenantEntity{}
		rsp.TotalCount = 0
		return nil
	}

	// 查询组织下所有商户
	ots, err := u.userStore.ListOrganizationExistTenants(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 根据名字查询商户
	tNames, err := u.userStore.SearchTenantIdsByName(ctx, req.GetTenantName())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	tempGols := make([]domain.OrganizationTenantIntf, 0)
	for _, o := range ots {
		if slices.Contains(tNames, o.GetTenantID()) {
			tempGols = append(tempGols, o)
		}
	}

	ots = tempGols

	queryIDs := make([]string, len(ots))
	for k, v := range ots {
		queryIDs[k] = v.GetTenantID()
	}

	// 获取商户方案
	treatments, err := u.userStore.GetTenantTreatments(ctx, queryIDs)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	tr_ids := make([]string, len(treatments))
	t_map := make(map[string]string)
	for k, each := range treatments {
		t_map[each.GetTenantID()] = each.GetTreatmentID()
		tr_ids[k] = each.GetTreatmentID()
	}
	tempGols = make([]domain.OrganizationTenantIntf, 0)
	if req.GetTreatmentId() == "null" {
		// 获取所有没有方案的商户
		for _, each := range ots {
			if _, ok := t_map[each.GetTenantID()]; !ok {
				tempGols = append(tempGols, each)
			}
		}
	} else if req.GetTreatmentId() != "" {
		// 获取指定方案的商户
		for _, each := range ots {
			if t, ok := t_map[each.GetTenantID()]; ok && t == req.GetTreatmentId() {
				tempGols = append(tempGols, each)
			}
		}
	} else if req.GetTreatmentId() == "" {
		tempGols = append(tempGols, ots...)
	}

	ots = tempGols

	// 查询商户订阅
	tids := make([]string, 0, len(ots))
	for _, each := range ots {
		tids = append(tids, each.GetTenantID())
	}
	ts, err := u.userStore.ListSubscriptionTimeline(ctx, tids)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 转换成商户订阅表
	ts_map := make(map[string]domain.SubscriptionTimelineIntf)
	pts_map := make(map[string]*pb.TenantSubscriptionTimeline)
	for _, each := range ts {
		ts_map[each.GetTenantID()] = each
		pts_map[each.GetTenantID()] = toProtoSubscriptionTimeline(each)
	}

	// 获取商品名称
	treatMap := make(map[string]string)
	if len(tr_ids) > 0 {
		listRsp, err := u.productAPI.ListTreatmentNameByIds(ctx, &productv1.ListTreatmentNameByIdsRequest{
			TreatmentIds: tr_ids,
		})
		if err != nil {
			return err
		}
		treatMap = listRsp.GetTreatmentNames()
	}

	treatNamesMap := make(map[string]string)
	for k, v := range t_map {
		treatNamesMap[k] = treatMap[v]
	}

	// 判断商户状态
	tenant_status_map := make(map[string]pb.TenantStatus)
	for _, each := range ots {
		if !each.GetIsActivated() {
			// 未认证
			tenant_status_map[each.GetTenantID()] = pb.TenantStatus_TENANT_STATUS_UNAUTH
		} else {
			if v, ok := ts_map[each.GetTenantID()]; !ok {
				// 没有订阅，待续期
				tenant_status_map[each.GetTenantID()] = pb.TenantStatus_TENANT_STATUS_PENDING
			} else {
				// 有订阅，判断是否过期
				tenant_status_map[each.GetTenantID()] = checkTenantSubscriptionTimeline(v)
			}
		}
	}

	goalOts := []domain.OrganizationTenantIntf{}
	// 查询 tenant
	// 查询状态
	if req.GetStatus() != pb.TenantStatus_TENANT_STATUS_UNSET {
		// 挑选出符合状态的 tenant
		for _, v := range ots {
			if tenant_status_map[v.GetTenantID()] == req.GetStatus() {
				goalOts = append(goalOts, v)
			}
		}
	} else {
		goalOts = ots
	}

	// 查询商户信息
	reviewFailTenant := []string{}
	tenantEntity := make([]*pb.TenantEntity, len(goalOts))
	for k, v := range goalOts {
		// 判断当前 tenant 的状态
		// 如果是审核状态则去获取(revision)，如果失败还要去获取失败原因
		// 其他情况获取 entity 即可
		if v.GetReviewStatus() != domain.TenantReviewStatusReviewSuccess {
			revision, err := u.userStore.GetLatestTenantEntityRevision(ctx, v.GetTenantID())
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
			if revision == nil {
				return errors.Errorf(ErrTenantNotFound, "revision[%s] not found", v.GetTenantID())
			}
			tenantEntity[k] = toProtoTenantEntityFromRevision(revision)
			tenantEntity[k].TenantStatus = tenant_status_map[v.GetTenantID()]
			tenantEntity[k].TenantReviewStatus = toProtoTenantReviewStatus(v.GetReviewStatus())
			if v.GetReviewStatus() == domain.TenantReviewStatusReviewFailed {
				// 获取失败原因
				reviewFailTenant = append(reviewFailTenant, v.GetTenantID())
			}
		} else {
			// 其他情况获取 tenantEntity 即可
			entity, err := u.userStore.GetTenantEntity(ctx, v.GetTenantID())
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}

			tenantEntity[k] = toProtoTenantEntity(entity)
			tenantEntity[k].TenantStatus = tenant_status_map[entity.GetTenantID()]
			tenantEntity[k].TenantReviewStatus = toProtoTenantReviewStatus(v.GetReviewStatus())
		}
	}

	// 获取失败原因
	if len(reviewFailTenant) != 0 {
		// 获取失败原因
		listRsp, err := u.reviewAPI.ListCertificateReviewStatus(ctx, &reviewpb.ListCertificateReviewStatusRequest{
			TenantIds: reviewFailTenant,
		})
		if err != nil {
			return err
		}
		result := listRsp.GetResults()
		for k, v := range tenantEntity {
			if res, ok := result[v.GetTenantId()]; ok {
				tenantEntity[k].FailReason = res.GetFailReason()
			}
		}
	}

	// 根据分页返回数据
	size := req.GetPagination().GetSize()
	offset := req.GetPagination().GetOffset()
	from := offset
	length := int32(len(tenantEntity))
	if from < length {
		res := tenantEntity[from:]
		if int32(len(res)) > size {
			res = tenantEntity[from : from+size]
		}
		rsp.Tenants = res
		rsp.Subscriptions = pts_map
		rsp.Treatments = treatNamesMap
	}

	rsp.TotalCount = length

	return nil
}

// 验证 request
func validateSearchTenantsRequest(req *pb.SearchTenantsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetStatus() == pb.TenantStatus_TENANT_STATUS_INVALID {
		return gerr.New("invalid SearchTenants status")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
