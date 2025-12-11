package domain

import "time"

func (t *TenantEntity) GetTenantID() string {
	if t == nil {
		return ""
	}
	return t.TenantID
}

func (t *TenantEntity) GetLogoUrl() string {
	if t == nil {
		return ""
	}
	return t.LogoUrl
}

func (t *TenantEntity) GetStoreName() string {
	if t == nil {
		return ""
	}
	return t.StoreName
}

func (t *TenantEntity) GetProvince() string {
	if t == nil {
		return ""
	}
	return t.Province
}

func (t *TenantEntity) GetCity() string {
	if t == nil {
		return ""
	}
	return t.City
}

func (t *TenantEntity) GetDistrict() string {
	if t == nil {
		return ""
	}
	return t.District
}

func (t *TenantEntity) GetStreet() string {
	if t == nil {
		return ""
	}
	return t.Street
}

func (t *TenantEntity) GetContactName() string {
	if t == nil {
		return ""
	}
	return t.ContactName
}

func (t *TenantEntity) GetContactPhone() string {
	if t == nil {
		return ""
	}
	return t.ContactPhone
}

func (t *TenantEntity) GetSocialCreditCode() string {
	if t == nil {
		return ""
	}
	return t.SocialCreditCode
}

func (t *TenantEntity) GetBusinessLicenseUrl() string {
	if t == nil {
		return ""
	}
	return t.BusinessLicenseUrl
}

func (t *TenantEntity) GetStaffCountQuota() int32 {
	if t == nil {
		return 0
	}
	return t.StaffCountQuota
}

func (t *TenantEntity) GetReportSharingStatus() bool {
	if t == nil {
		return false
	}
	return t.ReportSharingStatus
}

func (t *TenantEntity) GetConstitutionStatus() bool {
	if t == nil {
		return false
	}
	return t.ConstitutionStatus
}

func (t *TenantEntity) GetOverdue() int32 {
	if t == nil {
		return 0
	}
	return t.Overdue
}

func (t *TenantEntity) GetRev() int32 {
	if t == nil {
		return 0
	}
	return t.Rev
}

func (t *TenantEntity) GetSafePhone() string {
	if t == nil {
		return ""
	}
	return t.SafePhone
}

func (t *TenantEntity) GetCreatedAt() time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.CreatedAt
}

type TenantEntity struct {
	TenantID            string
	LogoUrl             string
	StoreName           string
	Province            string
	City                string
	District            string
	Street              string
	ContactName         string
	ContactPhone        string
	SocialCreditCode    string
	BusinessLicenseUrl  string
	StaffCountQuota     int32
	ReportSharingStatus bool
	ConstitutionStatus  bool
	Overdue             int32
	Rev                 int32
	SafePhone           string
	CreatedAt           time.Time
}

type TenantEntityMapper interface {
	ToDomainTenantEntityMapper
	FromDomainTenantEntityMapper
}

type ToDomainTenantEntityMapper interface {
	ToDomainTenantEntity() TenantEntityIntf
}

type FromDomainTenantEntityMapper interface {
	FromDomainTenantEntity(TenantEntityIntf)
}

type TenantEntityIntf interface {
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
	GetStaffCountQuota() int32
	GetReportSharingStatus() bool
	GetConstitutionStatus() bool
	GetOverdue() int32
	GetRev() int32
	GetSafePhone() string
	GetCreatedAt() time.Time
}

var _ TenantEntityIntf = (*TenantEntity)(nil)
