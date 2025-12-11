package review

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ReviewAPIHandler) ListReviewIssues(ctx context.Context, req *pb.ListReviewIssuesRequest, rsp *pb.ListReviewIssuesResponse) error {
	// 1.验证request
	err := validateListReviewIssuesRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询组织
	var reqOid []string
	if req.GetOrganizationName() != "" {
		getRsp, err := s.userAPI.GetOrganizationByName(ctx, &userpb.GetOrganizationByNameRequest{
			OrganizationName: req.GetOrganizationName(),
		})
		if err != nil {
			return err
		}
		// 如果查不到直接返回空即可
		if len(getRsp.GetOrganizations()) == 0 {
			return nil
		}
		reqOid = make([]string, len(getRsp.GetOrganizations()))
		for k, v := range getRsp.GetOrganizations() {
			reqOid[k] = v.GetOrganizationId()
		}
	}

	// 获取可查询组织权限
	lirRsp, err := s.userAPI.ListSystemUserDataPrivileges(ctx, &userpb.ListSystemUserDataPrivilegesRequest{
		UserId: req.GetReviewerUserId(),
	})
	if err != nil {
		return err
	}

	// 查询组织 && 可查询组织权限 的交集就是能查询到的
	oid := intersectionOrganizationIds(reqOid, lirRsp.GetOrganizationIds())
	if req.GetOrganizationName() == "" {
		oid = lirRsp.GetOrganizationIds()
	}

	pagination := req.GetPagination()
	// 如果要查询全部就去查询全部
	if req.GetSearchAll() {
		var issues []domain.ReviewIssueWithResultIntf
		var totalCount int64
		var err error
		if req.GetOrganizationName() != "" {
			// 查找OrganizationId的工单
			issues, totalCount, err = s.reviewStore.ListReviewIssuesWithOrganizationIds(ctx, int(pagination.GetSize()), int(pagination.GetOffset()), oid)
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
		} else {
			// 查找所有的工单
			issues, totalCount, err = s.reviewStore.ListReviewIssuesWithOrganizationIds(ctx, int(pagination.GetSize()), int(pagination.GetOffset()), oid)
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
		}
		// 数据转化
		pIssues, err := s.convertToAppIssues(ctx, issues)
		if err != nil {
			return err
		}
		rsp.TotalCount = int32(totalCount)
		rsp.ReviewIssues = pIssues
	} else {
		var issues []domain.ReviewIssueWithResultIntf
		var totalCount int64
		var err error
		status, err := toDomainReviewStatus(req.GetStatus())
		if err != nil {
			return errors.Error(codes.InvalidRequest, err.Error())
		}

		if req.GetOrganizationName() != "" {
			// 查找OrganizationId的工单
			issues, totalCount, err = s.reviewStore.ListReviewIssuesWithStatusAndOrganizationIds(ctx, status, int(pagination.GetSize()), int(pagination.GetOffset()), oid)
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
		} else {
			// 查找所有的工单
			issues, totalCount, err = s.reviewStore.ListReviewIssuesWithStatusAndOrganizationIds(ctx, status, int(pagination.GetSize()), int(pagination.GetOffset()), oid)
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
		}
		// 数据转化
		pIssues, err := s.convertToAppIssues(ctx, issues)
		if err != nil {
			return err
		}
		rsp.TotalCount = int32(totalCount)
		rsp.ReviewIssues = pIssues
	}
	return nil
}

// 验证request
func validateListReviewIssuesRequest(req *pb.ListReviewIssuesRequest) error {
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}

// listSystemUsers
func (u *ReviewAPIHandler) listSystemUsers(ctx context.Context, userIds []string) (map[string]*userpb.SystemUser, error) {
	book := make(map[string]bool)
	var ids []string
	for _, v := range userIds {
		if !book[v] {
			ids = append(ids, v)
		}
	}
	// 查询接单人姓名
	listRsp, err := u.userAPI.ListSystemUsersByIds(ctx, &userpb.ListSystemUsersByIdsRequest{
		SystemUserIds: ids,
	})
	if err != nil {
		return nil, err
	}
	return listRsp.GetSystemUsers(), nil
}

// convertToAppIssues
func (u *ReviewAPIHandler) convertToAppIssues(ctx context.Context, issues []domain.ReviewIssueWithResultIntf) ([]*pb.ReviewIssue, error) {
	if len(issues) == 0 {
		return []*pb.ReviewIssue{}, nil
	}
	// 获取所有工单接单员ID
	var userIds []string
	var oids []string
	var err error
	book := make(map[string]bool)
	for _, v := range issues {
		if v.GetReviewerUserID() != "" {
			// 如果不存在才保存
			if _, ok := book[v.GetReviewerUserID()]; !ok {
				userIds = append(userIds, v.GetReviewerUserID())
				book[v.GetReviewerUserID()] = true
			}
		}
		if v.GetSubmitterOrganizationID() != "" {
			// 如果不存在才保存
			if _, ok := book[v.GetSubmitterOrganizationID()]; !ok {
				oids = append(oids, v.GetSubmitterOrganizationID())
				book[v.GetSubmitterOrganizationID()] = true
			}
		}
	}
	userMap := make(map[string]*userpb.SystemUser)
	organizationNames := make(map[string]*userpb.Organization)
	if len(userIds) != 0 {
		// 获取所有接单者信息
		userMap, err = u.listSystemUsers(ctx, userIds)
		if err != nil {
			return nil, err
		}
	}
	if len(oids) != 0 {
		getRsp, err := u.userAPI.GetOrganizationsName(ctx, &userpb.GetOrganizationsNameRequest{
			OrganizationIds: oids,
		})
		if err != nil {
			return nil, err
		}
		organizationNames = getRsp.GetOrganizationName()
	}
	pIssues := make([]*pb.ReviewIssue, len(issues))
	// 构建工单信息返回数据
	for k, v := range issues {
		username := ""
		organizationName := ""
		if sysUser, ok := userMap[v.GetReviewerUserID()]; ok {
			username = sysUser.GetNickname()
		}
		if oname, ok := organizationNames[v.GetSubmitterOrganizationID()]; ok {
			organizationName = oname.GetName()
		}
		pIssues[k] = toProtoReviewIssue(v, username, organizationName)
	}
	return pIssues, nil
}

// 取组织id交集
func intersectionOrganizationIds(oid, pid []string) []string {
	if len(oid) == 0 || len(pid) == 0 {
		return []string{}
	}
	oidMap := make(map[string]bool)
	for _, v := range oid {
		oidMap[v] = true
	}
	// 目标id
	var goalId []string
	// 查看pid中和oid都包含的
	for _, v := range pid {
		if oidMap[v] {
			goalId = append(goalId, v)
		}
	}
	return goalId
}
