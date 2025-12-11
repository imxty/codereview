package review

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"gorm.io/gorm"
)

type TreatmentReviewResult struct {
	TreatmentRevID string          `gorm:"column:target_rev_id"`
	ReviewResultID string          `gorm:"primary_key;column:review_result_id"`
	ReviewIssueID  string          `gorm:"column:review_issue_id"`
	ReviewTime     time.Time       `gorm:"column:review_time"`
	ReviewUserID   string          `gorm:"column:review_user_id"`
	Result         bool            `gorm:"column:result"`
	Comment        string          `gorm:"column:comment"`
	Rev            int32           `gorm:"column:rev"`
	CreatedAt      time.Time       // 创建时间
	UpdatedAt      time.Time       // 更新时间
	DeletedAt      *gorm.DeletedAt // 删除时间
}

func (r TreatmentReviewResult) TableName() string {
	return "treatment_review_result"
}

// domain->db
func (r *TreatmentReviewResult) FromDomainTreatmentReviewResult(d domain.TreatmentReviewResultIntf) {
	if r == nil || d == nil {
		return
	}
	r.TreatmentRevID = d.GetTreatmentRevID()
	r.ReviewResultID = d.GetReviewResultID()
	r.ReviewIssueID = d.GetReviewIssueID()
	r.ReviewTime = d.GetReviewTime()
	r.ReviewUserID = d.GetReviewUserID()
	r.Result = d.GetResult()
	r.Comment = d.GetComment()
	r.Rev = d.GetRev()
}

// db->domain
func (r *TreatmentReviewResult) ToDomainTreatmentReviewResult() domain.TreatmentReviewResultIntf {
	if r == nil {
		return nil
	}

	p := domain.TreatmentReviewResult{
		TreatmentRevID: r.TreatmentRevID,
		ReviewResultID: r.ReviewResultID,
		ReviewIssueID:  r.ReviewIssueID,
		ReviewTime:     r.ReviewTime,
		ReviewUserID:   r.ReviewUserID,
		Result:         r.Result,
		Comment:        r.Comment,
		Rev:            r.Rev,
	}
	return &p
}
