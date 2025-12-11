package domain

import (
	"time"
)

type TreatmentReviewResult struct {
	TreatmentRevID string
	ReviewResultID string
	ReviewIssueID  string
	ReviewTime     time.Time
	ReviewUserID   string
	Result         bool
	Comment        string
	Rev            int32
}

type TreatmentReviewResultMapper interface {
	ToDomainTreatmentReviewResultMapper
	FromDomainTreatmentReviewResultMapper
}

type ToDomainTreatmentReviewResultMapper interface {
	ToDomainReviewResult() ReviewResultIntf
}

type FromDomainTreatmentReviewResultMapper interface {
	FromDomainReviewResult(ReviewResultIntf)
}

type TreatmentReviewResultIntf interface {
	GetTreatmentRevID() string
	GetReviewResultID() string
	GetReviewIssueID() string
	GetReviewTime() time.Time
	GetReviewUserID() string
	GetResult() bool
	GetComment() string
	GetRev() int32
}

var _ TreatmentReviewResultIntf = (*TreatmentReviewResult)(nil)

func (r *TreatmentReviewResult) GetTreatmentRevID() string {
	if r == nil {
		return ""
	}
	return r.TreatmentRevID
}
func (r *TreatmentReviewResult) GetReviewResultID() string {
	if r == nil {
		return ""
	}
	return r.ReviewResultID
}

func (r *TreatmentReviewResult) GetReviewIssueID() string {
	if r == nil {
		return ""
	}
	return r.ReviewIssueID
}

func (r *TreatmentReviewResult) GetReviewTime() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.ReviewTime
}

func (r *TreatmentReviewResult) GetReviewUserID() string {
	if r == nil {
		return ""
	}
	return r.ReviewUserID
}

func (r *TreatmentReviewResult) GetResult() bool {
	if r == nil {
		return false
	}
	return r.Result
}

func (r *TreatmentReviewResult) GetComment() string {
	if r == nil {
		return ""
	}
	return r.Comment
}

func (r *TreatmentReviewResult) GetRev() int32 {
	if r == nil {
		return 0
	}
	return r.Rev
}
