package domain

import "time"

const (
	// 审核中
	TenantReviewStatusReviewing = 0
	// 审核成功
	TenantReviewStatusReviewSuccess = 1
	// 审核失败
	TenantReviewStatusReviewFailed = 2
)

const (
	// 未认证
	TenantStatusUnAuth = 0
	// 待续期
	TenantStatusPending = 1
	// 使用中
	TenantStatusUsing = 2
)

type OrganizationTenant struct {
	OrganizationTenantID string
	OrganizationID       string
	TenantID             string
	IsActivated          bool
	ReviewStatus         int32
	Rev                  int32
	CreatedAt            time.Time
}

type OrganizationTenantMapper interface {
	ToDomainOrganizationTenantMapper
	FromDomainOrganizationTenantMapper
}

type ToDomainOrganizationTenantMapper interface {
	ToDomainOrganizationTenant() OrganizationTenantIntf
}

type FromDomainOrganizationTenantMapper interface {
	FromDomainOrganizationTenant(OrganizationTenantIntf)
}

type OrganizationTenantIntf interface {
	GetOrganizationTenantID() string
	GetOrganizationID() string
	GetTenantID() string
	GetIsActivated() bool
	GetReviewStatus() int32
	GetRev() int32
	GetCreatedAt() time.Time
}

var _ OrganizationTenantIntf = (*OrganizationTenant)(nil)

func (o *OrganizationTenant) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *OrganizationTenant) GetOrganizationTenantID() string {
	if o == nil {
		return ""
	}
	return o.OrganizationTenantID
}

func (o *OrganizationTenant) GetOrganizationID() string {
	if o == nil {
		return ""
	}
	return o.OrganizationID
}

func (o *OrganizationTenant) GetTenantID() string {
	if o == nil {
		return ""
	}
	return o.TenantID
}

func (o *OrganizationTenant) GetReviewStatus() int32 {
	if o == nil {
		return 0
	}
	return o.ReviewStatus
}

func (o *OrganizationTenant) GetRev() int32 {
	if o == nil {
		return 0
	}
	return o.Rev
}

func (o *OrganizationTenant) GetIsActivated() bool {
	if o == nil {
		return false
	}
	return o.IsActivated
}
