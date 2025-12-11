package domain

import "time"

func (t *TenantEntityRevision) GetTenantEntityRevisionID() string {
	if t == nil {
		return ""
	}
	return t.TenantEntityRevisionID
}

func (t *TenantEntityRevision) GetTenantID() string {
	if t == nil {
		return ""
	}
	return t.TenantID
}

func (t *TenantEntityRevision) GetLogoUrl() string {
	if t == nil {
		return ""
	}
	return t.LogoUrl
}

func (t *TenantEntityRevision) GetStoreName() string {
	if t == nil {
		return ""
	}
	return t.StoreName
}

func (t *TenantEntityRevision) GetProvince() string {
	if t == nil {
		return ""
	}
	return t.Province
}

func (t *TenantEntityRevision) GetCity() string {
	if t == nil {
		return ""
	}
	return t.City
}

func (t *TenantEntityRevision) GetDistrict() string {
	if t == nil {
		return ""
	}
	return t.District
}

func (t *TenantEntityRevision) GetStreet() string {
	if t == nil {
		return ""
	}
	return t.Street
}

func (t *TenantEntityRevision) GetContactName() string {
	if t == nil {
		return ""
	}
	return t.ContactName
}

func (t *TenantEntityRevision) GetContactPhone() string {
	if t == nil {
		return ""
	}
	return t.ContactPhone
}

func (t *TenantEntityRevision) GetSocialCreditCode() string {
	if t == nil {
		return ""
	}
	return t.SocialCreditCode
}

func (t *TenantEntityRevision) GetBusinessLicenseUrl() string {
	if t == nil {
		return ""
	}
	return t.BusinessLicenseUrl
}

func (t *TenantEntityRevision) GetRev() int32 {
	if t == nil {
		return 0
	}
	return t.Rev
}

func (t *TenantEntityRevision) GetFailReason() string {
	if t == nil {
		return ""
	}
	return t.FailReason
}

func (t *TenantEntityRevision) GetSafePhone() string {
	if t == nil {
		return ""
	}
	return t.SafePhone
}

func (t *TenantEntityRevision) GetCreatedAt() time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.CreatedAt
}

type TenantEntityRevision struct {
	TenantEntityRevisionID string
	TenantID               string
	LogoUrl                string
	StoreName              string
	Province               string
	City                   string
	District               string
	Street                 string
	ContactName            string
	ContactPhone           string
	SocialCreditCode       string
	BusinessLicenseUrl     string
	Rev                    int32
	FailReason             string
	SafePhone              string
	CreatedAt              time.Time
}

type TenantEntityRevisionMapper interface {
	ToDomainTenantEntityRevisionMapper
	FromDomainTenantEntityRevisionMapper
}

type ToDomainTenantEntityRevisionMapper interface {
	ToDomainTenantEntityRevision() TenantEntityRevisionIntf
}

type FromDomainTenantEntityRevisionMapper interface {
	FromDomainTenantEntityRevision(TenantEntityRevisionIntf)
}

type TenantEntityRevisionIntf interface {
	GetTenantEntityRevisionID() string
	GetTenantID() string
	GetLogoUrl() string
	GetStoreName() string
	GetProvince() string
	GetCity() string
	GetDistrict() string
	GetStreet() string
	GetContactName() string
	GetContactPhone() string
	GetSocialCreditCode() string
	GetBusinessLicenseUrl() string
	GetRev() int32
	GetFailReason() string
	GetSafePhone() string
	GetCreatedAt() time.Time
}

var _ TenantEntityRevisionIntf = (*TenantEntityRevision)(nil)
