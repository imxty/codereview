package review

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"gorm.io/gorm"
)

type ReviewIssue struct {
	ReviewIssueID           string          `gorm:"primary_key;column:review_issue_id"`
	SubmitterOrganizationID string          `gorm:"column:submitter_organization_id"`
	SubmitterTenantID       string          `gorm:"column:submitter_tenant_id"`
	SubmitTime              time.Time       `gorm:"column:submit_time"`
	ReviewerUserID          string          `gorm:"column:reviewer_user_id"`
	AcceptanceTime          time.Time       `gorm:"column:acceptance_time"`
	ReviewStatus            int32           `gorm:"column:review_status"`
	TargetType              int32           `gorm:"column:target_type"`
	IsCancel                bool            `gorm:"column:is_cancel"`
	TargetRevID             string          `gorm:"column:target_rev_id"`
	Rev                     int32           `gorm:"column:rev"`
	CreatedAt               time.Time       // 创建时间
	UpdatedAt               time.Time       // 更新时间
	DeletedAt               *gorm.DeletedAt // 删除时间
}

// 携带结果的工单
type ReviewIssueWithResult struct {
	ReviewIssueID           string `gorm:"primary_key;column:review_issue_id"`
	SubmitterOrganizationID string `gorm:"column:submitter_organization_id"`
	SubmitterTenantID       string `gorm:"column:submitter_tenant_id"`
	SubmitTime              time.Time
	ReviewerUserID          string `gorm:"column:reviewer_user_id"`
	AcceptanceTime          time.Time
	ReviewStatus            int32
	TargetType              int32
	IsCancel                bool
	TargetRevID             string `gorm:"column:target_rev_id"`
	ReviewRev               int32  `gorm:"column:review_rev"`
	ReviewResultID          string `gorm:"primary_key;column:review_result_id"`
	ReviewTime              time.Time
	ReviewUserID            string `gorm:"column:review_user_id"`
	Result                  bool
	Comment                 string
	ResultRev               int32           `gorm:"column:result_rev"`
	CreatedAt               time.Time       // 创建时间
	UpdatedAt               time.Time       // 更新时间
	DeletedAt               *gorm.DeletedAt // 删除时间
}

func (r ReviewIssue) TableName() string {
	return "review_issue"
}

// domain->db
func (r *ReviewIssue) FromDomainReviewIssue(d domain.ReviewIssueIntf) {
	if r == nil || d == nil {
		return
	}

	r.ReviewIssueID = d.GetReviewIssueID()
	r.SubmitterOrganizationID = d.GetSubmitterOrganizationID()
	r.SubmitterTenantID = d.GetSubmitterTenantID()
	r.SubmitTime = d.GetSubmitTime()
	r.ReviewerUserID = d.GetReviewerUserID()
	r.AcceptanceTime = d.GetAcceptanceTime()
	r.ReviewStatus = d.GetReviewStatus()
	r.TargetType = d.GetTargetType()
	r.IsCancel = d.GetIsCancel()
	r.TargetRevID = d.GetTargetRevID()
	r.Rev = d.GetRev()
}

// db->domain
func (r *ReviewIssue) ToDomainReviewIssue() domain.ReviewIssueIntf {
	if r == nil {
		return nil
	}

	p := domain.ReviewIssue{
		ReviewIssueID:           r.ReviewIssueID,
		SubmitterOrganizationID: r.SubmitterOrganizationID,
		SubmitterTenantID:       r.SubmitterTenantID,
		SubmitTime:              r.SubmitTime,
		ReviewerUserID:          r.ReviewerUserID,
		AcceptanceTime:          r.AcceptanceTime,
		ReviewStatus:            r.ReviewStatus,
		TargetType:              r.TargetType,
		IsCancel:                r.IsCancel,
		TargetRevID:             r.TargetRevID,
		Rev:                     r.Rev,
	}
	return &p
}

// domain->db
func (r *ReviewIssueWithResult) FromDomainReviewIssueWithResult(d domain.ReviewIssueWithResultIntf) {
	if r == nil || d == nil {
		return
	}

	r.ReviewIssueID = d.GetReviewIssueID()
	r.SubmitterOrganizationID = d.GetSubmitterOrganizationID()
	r.SubmitterTenantID = d.GetSubmitterTenantID()
	r.SubmitTime = d.GetSubmitTime()
	r.ReviewerUserID = d.GetReviewerUserID()
	r.AcceptanceTime = d.GetAcceptanceTime()
	r.ReviewStatus = d.GetReviewStatus()
	r.TargetType = d.GetTargetType()
	r.IsCancel = d.GetIsCancel()
	r.TargetRevID = d.GetTargetRevID()
	r.ReviewRev = d.GetReviewRev()
	r.ReviewResultID = d.GetReviewResultID()
	r.ReviewTime = d.GetReviewTime()
	r.ReviewUserID = d.GetReviewUserID()
	r.Result = d.GetResult()
	r.Comment = d.GetComment()
	r.ResultRev = d.GetResultRev()
}

// db->domain
func (r *ReviewIssueWithResult) ToDomainReviewIssueWithResult() domain.ReviewIssueWithResultIntf {
	if r == nil {
		return nil
	}

	p := domain.ReviewIssueWithResult{
		ReviewIssueID:           r.ReviewIssueID,
		SubmitterOrganizationID: r.SubmitterOrganizationID,
		SubmitterTenantID:       r.SubmitterTenantID,
		SubmitTime:              r.SubmitTime,
		ReviewerUserID:          r.ReviewerUserID,
		AcceptanceTime:          r.AcceptanceTime,
		ReviewStatus:            r.ReviewStatus,
		TargetType:              r.TargetType,
		IsCancel:                r.IsCancel,
		TargetRevID:             r.TargetRevID,
		ReviewRev:               r.ReviewRev,
		ReviewResultID:          r.ReviewResultID,
		ReviewTime:              r.ReviewTime,
		ReviewUserID:            r.ReviewUserID,
		Result:                  r.Result,
		Comment:                 r.Comment,
		ResultRev:               r.ResultRev,
	}
	return &p
}
