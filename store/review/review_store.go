package review

import (
	"context"
	"errors"
	"time"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	ptime "github.com/jinmukeji/huimaibao-service/pkg/time"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"gorm.io/gorm"
)

// ReviewStore 实现 Domain 层 ReviewStore 接口
type ReviewStore struct {
	// 数据库管理
	*dbutils.Connection
}

func NewReviewStore(management *dbutils.Connection) *ReviewStore {
	return &ReviewStore{
		management,
	}
}

// 实现 Domain 行为
var _ domain.ReviewRepository = (*ReviewStore)(nil)

// CreateReviewIssue
func (r *ReviewStore) CreateReviewIssue(ctx context.Context, issue domain.ReviewIssueIntf) error {
	var iss ReviewIssue
	iss.FromDomainReviewIssue(issue)
	return r.GetConnection(ctx).Model(&ReviewIssue{}).Create(&iss).Error
}

// CountMonthSubmitTreatmentReviewTimes
func (r *ReviewStore) CountMonthSubmitTreatmentReviewTimes(ctx context.Context, organizationId string) (int32, error) {
	var times int64
	now := time.Now().In(ptime.LocBeijing)
	startTime := ptime.GetFirstDateOfMonth(now)
	err := r.GetConnection(ctx).Model(&ReviewIssue{}).Where("submitter_organization_id = ? and target_type = ? and created_at between ? and ?", organizationId, domain.ReviewTypeProduct, startTime.UTC(), now.UTC()).
		Count(&times).Error
	if err != nil {
		return 0, err
	}
	return int32(times), nil
}

// ListOrganizationNotifications 获取组织通知
func (r *ReviewStore) ListOrganizationNotifications(ctx context.Context, organizationID string) ([]domain.ReviewNotificationIntf, error) {
	// 查询未读通知
	var notifications []ReviewNotification
	err := r.GetConnection(ctx).Model(&ReviewNotification{}).
		Where("organization_id = ? and has_read = 0", organizationID).
		Order("created_at desc").
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	dNotifications := make([]domain.ReviewNotificationIntf, len(notifications))
	for k, v := range notifications {
		dNotifications[k] = v.ToDomainReviewNotification()
	}
	return dNotifications, nil
}

// ListTenantNotifications 获取商户通知
func (r *ReviewStore) ListTenantNotifications(ctx context.Context, tenantID string) ([]domain.ReviewNotificationIntf, error) {
	// 查询未读通知
	var notifications []ReviewNotification
	err := r.GetConnection(ctx).Model(&ReviewNotification{}).
		Where("tenant_id = ? and has_read = 0", tenantID).
		Order("created_at desc").
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	dNotifications := make([]domain.ReviewNotificationIntf, len(notifications))
	for k, v := range notifications {
		dNotifications[k] = v.ToDomainReviewNotification()
	}
	return dNotifications, nil
}

// GetReviewIssueById 通过工单ID获取工单
func (r *ReviewStore) GetReviewIssueById(ctx context.Context, issueId string) (domain.ReviewIssueIntf, error) {
	var issue ReviewIssue
	err := r.GetConnection(ctx).Model(&ReviewIssue{}).Where("review_issue_id = ?", issueId).First(&issue).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return issue.ToDomainReviewIssue(), nil
}

// GetLatestReviewIssueByTargetId 通过目标ID获取最新工单
func (r *ReviewStore) GetLatestReviewIssueByTargetId(ctx context.Context, targetId string) (domain.ReviewIssueIntf, error) {
	var issue ReviewIssue
	err := r.GetConnection(ctx).Model(&ReviewIssue{}).Where("target_rev_id = ?", targetId).First(&issue).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return issue.ToDomainReviewIssue(), nil
}

// ListReviewIssueByTargetId 通过目标ID获取最新工单
func (r *ReviewStore) ListReviewIssueByTargetId(ctx context.Context, targetIds []string) ([]domain.ReviewIssueIntf, error) {
	var issues []ReviewIssue
	err := r.GetConnection(ctx).Model(&ReviewIssue{}).Where("target_rev_id IN (?)", targetIds).Find(&issues).Error
	if err != nil {
		return nil, err
	}
	dissues := make([]domain.ReviewIssueIntf, len(issues))
	for k, v := range issues {
		dissues[k] = v.ToDomainReviewIssue()
	}
	return dissues, nil
}

// AcceptReviewIssue 接单
func (r *ReviewStore) AcceptReviewIssue(ctx context.Context, issueId, uid string, rev int32) error {
	// 接单
	return r.GetConnection(ctx).Model(&ReviewIssue{}).Where("review_issue_id = ? and rev = ?", issueId, rev).Updates(map[string]interface{}{
		"review_status":    domain.ReviewStatusReviewing,
		"acceptance_time":  time.Now().UTC(),
		"reviewer_user_id": uid,
		"rev":              rev + 1,
	}).Error
}

// CancelReviewIssue 取消接单.
func (r *ReviewStore) CancelReviewIssue(ctx context.Context, issueId string, rev int32) error {
	// 取消接单
	return r.GetConnection(ctx).Model(&ReviewIssue{}).Where("review_issue_id = ? and rev = ?", issueId, rev).Updates(map[string]interface{}{
		"review_status":    domain.ReviewStatusPending,
		"acceptance_time":  ptime.InitTime,
		"reviewer_user_id": "",
		"rev":              rev + 1,
	}).Error
}

// CancelTreatmentReviewIssue 取消商品审核工单.
func (r *ReviewStore) CancelTreatmentReviewIssue(ctx context.Context, issueId string, rev int32) error {
	// 取消商品审核工单
	return r.GetConnection(ctx).Model(&ReviewIssue{}).Where("review_issue_id = ? and rev = ?", issueId, rev).Updates(map[string]interface{}{
		"is_cancel": 1,
		"rev":       rev + 1,
	}).Error
}

// GetReviewResultByIssueId 通过工单ID获取工单处理结果
func (r *ReviewStore) GetReviewResultByIssueId(ctx context.Context, issueId string) (domain.ReviewResultIntf, error) {
	var re ReviewResult
	err := r.GetConnection(ctx).Model(&ReviewResult{}).Where("review_issue_id = ?", issueId).First(&re).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return re.ToDomainReviewResult(), nil
}

// ListTenantCertificateIssues
func (r *ReviewStore) ListTenantCertificateIssues(ctx context.Context, tenantIds []string) ([]domain.ReviewIssueIntf, error) {
	var iss []ReviewIssue
	err := r.GetConnection(ctx).Model(&ReviewIssue{}).Raw(`SELECT 
        *
        FROM 
        review_issue 
        WHERE 
        review_issue_id IN (
            SELECT MAX(review_issue_id) 
            FROM 
            review_issue 
            WHERE 
            target_type = 0 AND submitter_tenant_id 
            IN 
            (?) 
            group by 
            submitter_tenant_id
        )`, tenantIds).Scan(&iss).Error
	if err != nil {
		return nil, err
	}
	dIss := make([]domain.ReviewIssueIntf, len(iss))
	for k, v := range iss {
		dIss[k] = v.ToDomainReviewIssue()
	}
	return dIss, nil
}

// ListReviewResultByIssueId 通过工单ID获取工单处理结果
func (r *ReviewStore) ListReviewResultByIssueId(ctx context.Context, issueId []string) ([]domain.ReviewResultIntf, error) {
	var re []ReviewResult
	err := r.GetConnection(ctx).Model(&ReviewResult{}).Where("review_issue_id IN (?)", issueId).Find(&re).Error
	if err != nil {
		return nil, err
	}
	dr := make([]domain.ReviewResultIntf, len(re))
	for k, v := range re {
		dr[k] = v.ToDomainReviewResult()
	}
	return dr, nil
}

// CommitReviewResult 提交审核结果
func (r *ReviewStore) CommitReviewResult(ctx context.Context, issueResult domain.ReviewResultIntf) error {
	var result ReviewResult
	result.FromDomainReviewResult(issueResult)
	return r.GetConnection(ctx).Model(&ReviewResult{}).Create(&result).Error
}

// CreateReviewNotification 创建审核通知.
func (r *ReviewStore) CreateReviewNotification(ctx context.Context, notification domain.ReviewNotificationIntf) error {
	var noti ReviewNotification
	noti.FromDomainReviewNotification(notification)
	return r.GetConnection(ctx).Model(&ReviewNotification{}).Create(&noti).Error
}

// PassReviewIssue 通过审核工单
func (r *ReviewStore) PassReviewIssue(ctx context.Context, issueId string, rev int32) error {
	return r.GetConnection(ctx).Model(&ReviewIssue{}).Where("review_issue_id = ? and rev = ?", issueId, rev).Updates(map[string]interface{}{
		"review_status": domain.ReviewStatusReviewSuccess,
		"rev":           rev + 1,
	}).Error
}

// ReviewIssueFailed 审核工单失败
func (r *ReviewStore) ReviewIssueFailed(ctx context.Context, issueId string, rev int32) error {
	return r.GetConnection(ctx).Model(&ReviewIssue{}).Where("review_issue_id = ? and rev = ?", issueId, rev).Updates(map[string]interface{}{
		"review_status": domain.ReviewStatusReviewFailed,
		"rev":           rev + 1,
	}).Error
}

// GetReviewNotification 获取审核通知.
func (r *ReviewStore) GetReviewNotification(ctx context.Context, notificationId string) (domain.ReviewNotificationIntf, error) {
	var noti ReviewNotification
	err := r.GetConnection(ctx).Model(&ReviewNotification{}).Where("review_notification_id = ?", notificationId).First(&noti).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return noti.ToDomainReviewNotification(), nil
}

// ConfirmReviewNotification 确认审核通知.
func (r *ReviewStore) ConfirmReviewNotification(ctx context.Context, notificationId string, rev int32) error {
	return r.GetConnection(ctx).Model(&ReviewNotification{}).
		Where("review_notification_id = ? and rev = ?", notificationId, rev).
		Updates(map[string]interface{}{
			"has_read": 1,
			"rev":      rev + 1,
		}).Error
}

// SearchAcceptableIssue 搜索可接单的工单
func (r *ReviewStore) SearchAcceptableIssue(ctx context.Context) (domain.ReviewIssueIntf, error) {
	var issue ReviewIssue
	err := r.GetConnection(ctx).Model(&ReviewIssue{}).Where("review_status = ?", domain.ReviewStatusPending).First(&issue).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return issue.ToDomainReviewIssue(), nil
}

// ListReviewIssues 获取工单列表
func (r *ReviewStore) ListReviewIssues(ctx context.Context, size, offset int) ([]domain.ReviewIssueWithResultIntf, int64, error) {
	var issues []ReviewIssueWithResult
	var count int64
	// 查询工单
	err := r.GetConnection(ctx).Raw(`SELECT 
        a.review_issue_id,
        a.submitter_tenant_id,
        a.submitter_organization_id,
        a.submit_time,
        a.reviewer_user_id,
        a.acceptance_time,
        a.review_status,
        a.target_type,
        a.is_cancel,
        a.target_rev_id,
        a.rev as review_rev,
        b.review_result_id,
        b.review_time,
        b.result,
        b.comment,
        b.rev as result_rev
    FROM review_issue a LEFT JOIN review_result b ON a.review_issue_id = b.review_issue_id
        WHERE
        b.deleted_at is NULL 
    ORDER BY
        a.is_cancel,
        CASE WHEN a.review_status >= 2 THEN 3
        WHEN a.review_status = 1 THEN 2
        ELSE 1 END,a.submit_time desc
        LIMIT ? OFFSET ?`, size, offset).Scan(&issues).Error
	if err != nil {
		return nil, 0, err
	}
	// 计算个数
	err = r.GetConnection(ctx).Model(&ReviewIssue{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	dIssues := make([]domain.ReviewIssueWithResultIntf, len(issues))
	for k, v := range issues {
		dIssues[k] = v.ToDomainReviewIssueWithResult()
	}
	return dIssues, count, nil
}

// ListReviewIssuesWithStatus 根据状态获取工单列表
func (r *ReviewStore) ListReviewIssuesWithStatus(ctx context.Context, status int32, size, offset int) ([]domain.ReviewIssueWithResultIntf, int64, error) {
	var issues []ReviewIssueWithResult
	var count int64
	// 查询工单
	err := r.GetConnection(ctx).Raw(`SELECT 
        a.review_issue_id,
        a.submitter_tenant_id,
        a.submitter_organization_id,
        a.submit_time,
        a.reviewer_user_id,
        a.acceptance_time,
        a.review_status,
        a.target_type,
        a.is_cancel,
        a.target_rev_id,
        a.rev as review_rev,
        b.review_result_id,
        b.review_time,
        b.reviewer_user_id,
        b.result,
        b.comment,
        b.rev as result_rev
        FROM review_issue a 
        LEFT JOIN review_result b ON a.review_issue_id = b.review_issue_id 
        WHERE
        a.review_status = ? and a.is_cancel = 0
        AND
        b.deleted_at is NULL
        ORDER BY
            a.is_cancel,
            CASE WHEN a.review_status >= 2 THEN 3
            WHEN a.review_status = 1 THEN 2
            ELSE 1 END,a.submit_time desc
        LIMIT ? OFFSET ?`, status, size, offset).Scan(&issues).Error
	if err != nil {
		return nil, 0, err
	}
	// 计算个数
	err = r.GetConnection(ctx).Model(&ReviewIssue{}).Where("review_status = ? and is_cancel = 0", status).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	dIssues := make([]domain.ReviewIssueWithResultIntf, len(issues))
	for k, v := range issues {
		dIssues[k] = v.ToDomainReviewIssueWithResult()
	}
	return dIssues, count, nil
}

// ListReviewIssuesWithOrganizationId 获取工单列表withOrganizationId
func (r *ReviewStore) ListReviewIssuesWithOrganizationIds(ctx context.Context, size, offset int, oids []string) ([]domain.ReviewIssueWithResultIntf, int64, error) {
	var issues []ReviewIssueWithResult
	var count int64
	// 查询工单
	err := r.GetConnection(ctx).Raw(`SELECT 
        a.review_issue_id,
        a.submitter_tenant_id,
        a.submitter_organization_id,
        a.submit_time,
        a.reviewer_user_id,
        a.acceptance_time,
        a.review_status,
        a.target_type,
        a.is_cancel,
        a.target_rev_id,
        a.rev as review_rev,
        b.review_result_id,
        b.review_time,
        b.result,
        b.comment,
        b.rev as result_rev
    FROM review_issue a LEFT JOIN review_result b ON a.review_issue_id = b.review_issue_id
        WHERE
        b.deleted_at is NULL and a.submitter_organization_id IN (?)
    ORDER BY
        a.is_cancel,
        CASE WHEN a.review_status >= 2 THEN 3
        WHEN a.review_status = 1 THEN 2
        ELSE 1 END,a.submit_time desc
        LIMIT ? OFFSET ?`, oids, size, offset).Scan(&issues).Error
	if err != nil {
		return nil, 0, err
	}
	// 计算个数
	err = r.GetConnection(ctx).Model(&ReviewIssue{}).Where("submitter_organization_id IN (?)", oids).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	dIssues := make([]domain.ReviewIssueWithResultIntf, len(issues))
	for k, v := range issues {
		dIssues[k] = v.ToDomainReviewIssueWithResult()
	}
	return dIssues, count, nil
}

// ListReviewIssuesWithStatusAndOrganizationId 根据状态获取工单列表withOrganizationId
func (r *ReviewStore) ListReviewIssuesWithStatusAndOrganizationIds(ctx context.Context, status int32, size, offset int, oids []string) ([]domain.ReviewIssueWithResultIntf, int64, error) {
	var issues []ReviewIssueWithResult
	var count int64
	// 查询工单
	err := r.GetConnection(ctx).Raw(`SELECT 
        a.review_issue_id,
        a.submitter_tenant_id,
        a.submitter_organization_id,
        a.submit_time,
        a.reviewer_user_id,
        a.acceptance_time,
        a.review_status,
        a.target_type,
        a.is_cancel,
        a.target_rev_id,
        a.rev as review_rev,
        b.review_result_id,
        b.review_time,
        b.reviewer_user_id,
        b.result,
        b.comment,
        b.rev as result_rev
        FROM review_issue a 
        LEFT JOIN review_result b ON a.review_issue_id = b.review_issue_id 
        WHERE
        a.review_status = ? and a.is_cancel = 0 and a.submitter_organization_id IN (?)
        AND
        b.deleted_at is NULL 
        ORDER BY
            a.is_cancel,
            CASE WHEN a.review_status >= 2 THEN 3
            WHEN a.review_status = 1 THEN 2
            ELSE 1 END,a.submit_time desc
        LIMIT ? OFFSET ?`, status, oids, size, offset).Scan(&issues).Error
	if err != nil {
		return nil, 0, err
	}
	// 计算个数
	err = r.GetConnection(ctx).Model(&ReviewIssue{}).Where("review_status = ? and is_cancel = 0 and submitter_organization_id IN (?)", status, oids).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	dIssues := make([]domain.ReviewIssueWithResultIntf, len(issues))
	for k, v := range issues {
		dIssues[k] = v.ToDomainReviewIssueWithResult()
	}
	return dIssues, count, nil
}

// CountReviewStatus 统计未审核，预付款未审核数量
func (r *ReviewStore) CountReviewStatus(ctx context.Context) (int32, int32, error) {
	var reviewCount int64
	var prestoreReviewCount int64
	err := r.GetConnection(ctx).Model(&ReviewIssue{}).Where("((review_status = ? or review_status = ?) and is_cancel = 0)", domain.ReviewStatusPending, domain.ReviewStatusReviewing).Count(&reviewCount).Error
	if err != nil {
		return 0, 0, err
	}
	err = r.GetConnection(ctx).Raw(`
        SELECT count(*) 
        FROM 
        tenant_prestore_review
        WHERE
        status = ?
    `, domain.ReviewStatusReviewing).Scan(&prestoreReviewCount).Error
	if err != nil {
		return 0, 0, err
	}
	return int32(reviewCount), int32(prestoreReviewCount), nil
}
