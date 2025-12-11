package review

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"gorm.io/gorm"
)

type ReviewResult struct {
	ReviewResultID string          `gorm:"primary_key;column:review_result_id"`
	ReviewIssueID  string          `gorm:"column:review_issue_id"`
	ReviewTime     time.Time       `gorm:"column:review_time"`
	ReviewerUserID string          `gorm:"column:reviewer_user_id"`
	Result         bool            `gorm:"column:result"`
	Comment        string          `gorm:"column:comment"`
	Rev            int32           `gorm:"column:rev"`
	CreatedAt      time.Time       // 创建时间
	UpdatedAt      time.Time       // 更新时间
	DeletedAt      *gorm.DeletedAt // 删除时间
}

func (r ReviewResult) TableName() string {
	return "review_result"
}

// domain->db
func (r *ReviewResult) FromDomainReviewResult(d domain.ReviewResultIntf) {
	if r == nil || d == nil {
		return
	}

	r.ReviewResultID = d.GetReviewResultID()
	r.ReviewIssueID = d.GetReviewIssueID()
	r.ReviewTime = d.GetReviewTime()
	r.ReviewerUserID = d.GetReviewerUserID()
	r.Result = d.GetResult()
	r.Comment = d.GetComment()
	r.Rev = d.GetRev()
}

// db->domain
func (r *ReviewResult) ToDomainReviewResult() domain.ReviewResultIntf {
	if r == nil {
		return nil
	}

	p := domain.ReviewResult{
		ReviewResultID: r.ReviewResultID,
		ReviewIssueID:  r.ReviewIssueID,
		ReviewTime:     r.ReviewTime,
		ReviewerUserID: r.ReviewerUserID,
		Result:         r.Result,
		Comment:        r.Comment,
		Rev:            r.Rev,
	}
	return &p
}
