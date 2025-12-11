package domain

import (
	"context"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
)

type ReviewRepository interface {
	dbutils.Tx

	// CreateReviewIssue
	CreateReviewIssue(ctx context.Context, issue ReviewIssueIntf) error
	// CountMonthSubmitTreatmentReviewTimes
	CountMonthSubmitTreatmentReviewTimes(ctx context.Context, organizationId string) (int32, error)
	// ListOrganizationNotifications 获取组织通知
	ListOrganizationNotifications(ctx context.Context, organizationID string) ([]ReviewNotificationIntf, error)
	// ListTenantNotifications 获取商户通知
	ListTenantNotifications(ctx context.Context, tenantID string) ([]ReviewNotificationIntf, error)

	// GetReviewIssueById 通过工单ID获取工单
	GetReviewIssueById(ctx context.Context, issueId string) (ReviewIssueIntf, error)
	// GetLatestReviewIssueByTargetId 通过目标ID获取最新工单
	GetLatestReviewIssueByTargetId(ctx context.Context, targetId string) (ReviewIssueIntf, error)
	// ListReviewIssueByTargetId 通过目标ID获取最新工单
	ListReviewIssueByTargetId(ctx context.Context, targetIds []string) ([]ReviewIssueIntf, error)
	// AcceptReviewIssue 接单
	AcceptReviewIssue(ctx context.Context, issueId, uid string, rev int32) error
	// CancelReviewIssue 取消接单.
	CancelReviewIssue(ctx context.Context, issueId string, rev int32) error
	// CancelTreatmentReviewIssue 取消商品审核工单.
	CancelTreatmentReviewIssue(ctx context.Context, issueId string, rev int32) error
	// GetReviewResultByIssueId 通过工单ID获取工单处理结果
	GetReviewResultByIssueId(ctx context.Context, issueId string) (ReviewResultIntf, error)
	// ListReviewResultByIssueId 通过工单ID获取工单处理结果
	ListReviewResultByIssueId(ctx context.Context, issueId []string) ([]ReviewResultIntf, error)
	// CommitReviewResult 提交审核结果
	CommitReviewResult(ctx context.Context, issueResult ReviewResultIntf) error
	// CreateReviewNotification 创建审核通知.
	CreateReviewNotification(ctx context.Context, notification ReviewNotificationIntf) error
	// PassReviewIssue 通过审核工单
	PassReviewIssue(ctx context.Context, issueId string, rev int32) error
	// ReviewIssueFailed 审核工单失败
	ReviewIssueFailed(ctx context.Context, issueId string, rev int32) error
	// ListTenantCertificateIssues
	ListTenantCertificateIssues(ctx context.Context, tenantIds []string) ([]ReviewIssueIntf, error)

	// GetReviewNotification 获取审核通知.
	GetReviewNotification(ctx context.Context, notificationId string) (ReviewNotificationIntf, error)
	// ConfirmReviewNotification 确认审核通知.
	ConfirmReviewNotification(ctx context.Context, notificationId string, rev int32) error

	// SearchAcceptableIssue 搜索可接单的工单
	SearchAcceptableIssue(ctx context.Context) (ReviewIssueIntf, error)
	// ListReviewIssues 获取工单列表
	ListReviewIssues(ctx context.Context, size, offset int) ([]ReviewIssueWithResultIntf, int64, error)
	// ListReviewIssuesWithStatus 根据状态获取工单列表
	ListReviewIssuesWithStatus(ctx context.Context, status int32, size, offset int) ([]ReviewIssueWithResultIntf, int64, error)

	// ListReviewIssuesWithOrganizationId 获取工单列表withOrganizationId
	ListReviewIssuesWithOrganizationIds(ctx context.Context, size, offset int, oids []string) ([]ReviewIssueWithResultIntf, int64, error)
	// ListReviewIssuesWithStatusAndOrganizationId 根据状态获取工单列表withOrganizationId
	ListReviewIssuesWithStatusAndOrganizationIds(ctx context.Context, status int32, size, offset int, oids []string) ([]ReviewIssueWithResultIntf, int64, error)

	// CountReviewStatus 统计未审核，预付款未审核数量
	CountReviewStatus(ctx context.Context) (int32, int32, error)
}
