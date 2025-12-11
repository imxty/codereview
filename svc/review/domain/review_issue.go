package domain

import (
	"time"
)

const (
	// 待审核
	ReviewStatusPending = 0
	// 审核中
	ReviewStatusReviewing = 1
	// 审核失败
	ReviewStatusReviewFailed = 2
	// 审核成功
	ReviewStatusReviewSuccess = 3
)

const (
	// 商户资质审核
	ReviewTypeCertificate = 0
	// 商品审核
	ReviewTypeProduct = 1
)

const (
	// 审核没有取消
	ReviewIsNotCanceled = 0
	// 审核取消
	ReviewIsCanceled = 1
)

type ReviewIssue struct {
	ReviewIssueID           string
	SubmitterOrganizationID string
	SubmitterTenantID       string
	SubmitTime              time.Time
	ReviewerUserID          string
	AcceptanceTime          time.Time
	ReviewStatus            int32
	TargetType              int32
	IsCancel                bool
	TargetRevID             string
	Rev                     int32
}

type ReviewIssueMapper interface {
	ToDomainReviewIssueMapper
	FromDomainReviewIssueMapper
}

type ToDomainReviewIssueMapper interface {
	ToDomainReviewIssue() ReviewIssueIntf
}

type FromDomainReviewIssueMapper interface {
	FromDomainReviewIssue(ReviewIssueIntf)
}

type ReviewIssueIntf interface {
	GetReviewIssueID() string
	GetSubmitterOrganizationID() string
	GetSubmitterTenantID() string
	GetSubmitTime() time.Time
	GetReviewerUserID() string
	GetAcceptanceTime() time.Time
	GetReviewStatus() int32
	GetTargetType() int32
	GetIsCancel() bool
	GetTargetRevID() string
	GetRev() int32
}

var _ ReviewIssueIntf = (*ReviewIssue)(nil)

type ReviewIssueWithResult struct {
	ReviewIssueID           string
	SubmitterOrganizationID string
	SubmitterTenantID       string
	SubmitTime              time.Time
	ReviewerUserID          string
	AcceptanceTime          time.Time
	ReviewStatus            int32
	TargetType              int32
	IsCancel                bool
	TargetRevID             string
	ReviewRev               int32
	ReviewResultID          string
	ReviewTime              time.Time
	ReviewUserID            string
	Result                  bool
	Comment                 string
	ResultRev               int32
}

type ReviewIssueWithResultMapper interface {
	ToDomainReviewIssueWithResultMapper
	FromDomainReviewIssueWithResultMapper
}

type ToDomainReviewIssueWithResultMapper interface {
	ToDomainReviewIssueWithResult() ReviewIssueWithResultIntf
}

type FromDomainReviewIssueWithResultMapper interface {
	FromDomainReviewIssueWithResult(ReviewIssueWithResultIntf)
}

type ReviewIssueWithResultIntf interface {
	GetReviewIssueID() string
	GetSubmitterOrganizationID() string
	GetSubmitterTenantID() string
	GetSubmitTime() time.Time
	GetReviewerUserID() string
	GetAcceptanceTime() time.Time
	GetReviewStatus() int32
	GetTargetType() int32
	GetIsCancel() bool
	GetTargetRevID() string
	GetReviewRev() int32
	GetReviewResultID() string
	GetReviewTime() time.Time
	GetReviewUserID() string
	GetResult() bool
	GetComment() string
	GetResultRev() int32
}

var _ ReviewIssueWithResultIntf = (*ReviewIssueWithResult)(nil)

func (r *ReviewIssue) GetReviewIssueID() string {
	if r == nil {
		return ""
	}
	return r.ReviewIssueID
}

func (r *ReviewIssue) GetSubmitterOrganizationID() string {
	if r == nil {
		return ""
	}
	return r.SubmitterOrganizationID
}

func (r *ReviewIssue) GetSubmitterTenantID() string {
	if r == nil {
		return ""
	}
	return r.SubmitterTenantID
}

func (r *ReviewIssue) GetSubmitTime() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.SubmitTime
}

func (r *ReviewIssue) GetReviewerUserID() string {
	if r == nil {
		return ""
	}
	return r.ReviewerUserID
}

func (r *ReviewIssue) GetAcceptanceTime() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.AcceptanceTime
}

func (r *ReviewIssue) GetReviewStatus() int32 {
	if r == nil {
		return 0
	}
	return r.ReviewStatus
}

func (r *ReviewIssue) GetTargetType() int32 {
	if r == nil {
		return 0
	}
	return r.TargetType
}

func (r *ReviewIssue) GetIsCancel() bool {
	if r == nil {
		return false
	}
	return r.IsCancel
}

func (r *ReviewIssue) GetTargetRevID() string {
	if r == nil {
		return ""
	}
	return r.TargetRevID
}

func (r *ReviewIssue) GetRev() int32 {
	if r == nil {
		return 0
	}
	return r.Rev
}

func (r *ReviewIssueWithResult) GetReviewIssueID() string {
	if r == nil {
		return ""
	}
	return r.ReviewIssueID
}

func (r *ReviewIssueWithResult) GetSubmitterOrganizationID() string {
	if r == nil {
		return ""
	}
	return r.SubmitterOrganizationID
}

func (r *ReviewIssueWithResult) GetSubmitterTenantID() string {
	if r == nil {
		return ""
	}
	return r.SubmitterTenantID
}

func (r *ReviewIssueWithResult) GetSubmitTime() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.SubmitTime
}

func (r *ReviewIssueWithResult) GetReviewerUserID() string {
	if r == nil {
		return ""
	}
	return r.ReviewerUserID
}

func (r *ReviewIssueWithResult) GetAcceptanceTime() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.AcceptanceTime
}

func (r *ReviewIssueWithResult) GetReviewStatus() int32 {
	if r == nil {
		return 0
	}
	return r.ReviewStatus
}

func (r *ReviewIssueWithResult) GetTargetType() int32 {
	if r == nil {
		return 0
	}
	return r.TargetType
}

func (r *ReviewIssueWithResult) GetIsCancel() bool {
	if r == nil {
		return false
	}
	return r.IsCancel
}

func (r *ReviewIssueWithResult) GetTargetRevID() string {
	if r == nil {
		return ""
	}
	return r.TargetRevID
}

func (r *ReviewIssueWithResult) GetReviewRev() int32 {
	if r == nil {
		return 0
	}
	return r.ReviewRev
}

func (r *ReviewIssueWithResult) GetReviewResultID() string {
	if r == nil {
		return ""
	}
	return r.ReviewResultID
}

func (r *ReviewIssueWithResult) GetReviewTime() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.ReviewTime
}

func (r *ReviewIssueWithResult) GetReviewUserID() string {
	if r == nil {
		return ""
	}
	return r.ReviewUserID
}

func (r *ReviewIssueWithResult) GetResult() bool {
	if r == nil {
		return false
	}
	return r.Result
}

func (r *ReviewIssueWithResult) GetComment() string {
	if r == nil {
		return ""
	}
	return r.Comment
}

func (r *ReviewIssueWithResult) GetResultRev() int32 {
	if r == nil {
		return 0
	}
	return r.ResultRev
}
