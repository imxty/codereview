package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

type TenantEntityRevision struct {
	TenantEntityRevisionID string          `gorm:"primary_key;column:tenant_entity_revision_id"`
	TenantID               string          `gorm:"column:tenant_id"`
	LogoUrl                string          `gorm:"column:logo_url"`
	SafePhone              string          `gorm:"column:safe_phone"`
	StoreName              string          `gorm:"column:store_name"`
	Province               string          `gorm:"column:province"`
	City                   string          `gorm:"column:city"`
	District               string          `gorm:"column:district"`
	Street                 string          `gorm:"column:street"`
	ContactName            string          `gorm:"column:contact_name"`
	ContactPhone           string          `gorm:"column:contact_phone"`
	SocialCreditCode       string          `gorm:"column:social_credit_code"`
	BusinessLicenseUrl     string          `gorm:"column:business_license_url"`
	Rev                    int32           `gorm:"column:rev"`
	FailReason             string          `gorm:"column:fail_reason"`
	CreatedAt              time.Time       // 创建时间
	UpdatedAt              time.Time       // 更新时间
	DeletedAt              *gorm.DeletedAt // 删除时间
}

func (t TenantEntityRevision) TableName() string {
	return "tenant_entity_revision"
}

// domain->db
func (t *TenantEntityRevision) FromDomainTenantEntityRevision(d domain.TenantEntityRevisionIntf) {
	if t == nil || d == nil {
		return
	}

	t.TenantEntityRevisionID = d.GetTenantEntityRevisionID()
	t.TenantID = d.GetTenantID()
	t.LogoUrl = d.GetLogoUrl()
	t.SafePhone = d.GetSafePhone()
	t.StoreName = d.GetStoreName()
	t.Province = d.GetProvince()
	t.City = d.GetCity()
	t.District = d.GetDistrict()
	t.Street = d.GetStreet()
	t.ContactName = d.GetContactName()
	t.ContactPhone = d.GetContactPhone()
	t.SocialCreditCode = d.GetSocialCreditCode()
	t.BusinessLicenseUrl = d.GetBusinessLicenseUrl()
	t.Rev = d.GetRev()
	t.FailReason = d.GetFailReason()
	t.SafePhone = d.GetSafePhone()
}

// db->domain
func (t *TenantEntityRevision) ToDomainTenantEntityRevision() domain.TenantEntityRevisionIntf {
	if t == nil {
		return nil
	}

	p := domain.TenantEntityRevision{
		TenantEntityRevisionID: t.TenantEntityRevisionID,
		TenantID:               t.TenantID,
		LogoUrl:                t.LogoUrl,
		StoreName:              t.StoreName,
		Province:               t.Province,
		City:                   t.City,
		District:               t.District,
		Street:                 t.Street,
		ContactName:            t.ContactName,
		ContactPhone:           t.ContactPhone,
		SocialCreditCode:       t.SocialCreditCode,
		BusinessLicenseUrl:     t.BusinessLicenseUrl,
		Rev:                    t.Rev,
		FailReason:             t.FailReason,
		SafePhone:              t.SafePhone,
	}
	return &p
}
