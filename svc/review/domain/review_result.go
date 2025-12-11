package domain

import (
	"time"
)

const (
	// 未通过
	ResultFailed = 0
	// 通过
	ResultPassed = 1
)

type ReviewResult struct {
	ReviewResultID string
	ReviewIssueID  string
	ReviewTime     time.Time
	ReviewerUserID string
	Result         bool
	Comment        string
	Rev            int32
}

type ReviewResultMapper interface {
	ToDomainReviewResultMapper
	FromDomainReviewResultMapper
}

type ToDomainReviewResultMapper interface {
	ToDomainReviewResult() ReviewResultIntf
}

type FromDomainReviewResultMapper interface {
	FromDomainReviewResult(ReviewResultIntf)
}

type ReviewResultIntf interface {
	GetReviewResultID() string
	GetReviewIssueID() string
	GetReviewTime() time.Time
	GetReviewerUserID() string
	GetResult() bool
	GetComment() string
	GetRev() int32
}

var _ ReviewResultIntf = (*ReviewResult)(nil)

func (r *ReviewResult) GetReviewResultID() string {
	if r == nil {
		return ""
	}
	return r.ReviewResultID
}

func (r *ReviewResult) GetReviewIssueID() string {
	if r == nil {
		return ""
	}
	return r.ReviewIssueID
}

func (r *ReviewResult) GetReviewTime() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.ReviewTime
}

func (r *ReviewResult) GetReviewerUserID() string {
	if r == nil {
		return ""
	}
	return r.ReviewerUserID
}

func (r *ReviewResult) GetResult() bool {
	if r == nil {
		return false
	}
	return r.Result
}

func (r *ReviewResult) GetComment() string {
	if r == nil {
		return ""
	}
	return r.Comment
}

func (r *ReviewResult) GetRev() int32 {
	if r == nil {
		return 0
	}
	return r.Rev
}
