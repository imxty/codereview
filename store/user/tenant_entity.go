package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

type TenantEntity struct {
	TenantID            string          `gorm:"primary_key;column:tenant_id"`
	LogoUrl             string          `gorm:"column:logo_url"`
	SafePhone           string          `gorm:"column:safe_phone"`
	StoreName           string          `gorm:"column:store_name"`
	Province            string          `gorm:"column:province"`
	City                string          `gorm:"column:city"`
	District            string          `gorm:"column:district"`
	Street              string          `gorm:"column:street"`
	ContactName         string          `gorm:"column:contact_name"`
	ContactPhone        string          `gorm:"column:contact_phone"`
	SocialCreditCode    string          `gorm:"column:social_credit_code"`
	BusinessLicenseUrl  string          `gorm:"column:business_license_url"`
	StaffCountQuota     int32           `gorm:"column:staff_count_quota"`
	ReportSharingStatus bool            `gorm:"column:report_sharing_status"`
	ConstitutionStatus  bool            `gorm:"column:constitution_status"`
	Overdue             int32           `gorm:"column:overdue"`
	Rev                 int32           `gorm:"column:rev"`
	CreatedAt           time.Time       // 创建时间
	UpdatedAt           time.Time       // 更新时间
	DeletedAt           *gorm.DeletedAt // 删除时间
}

func (t TenantEntity) TableName() string {
	return "tenant_entity"
}

// domain->db
func (t *TenantEntity) FromDomainTenantEntity(d domain.TenantEntityIntf) {
	if t == nil || d == nil {
		return
	}

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
	t.StaffCountQuota = d.GetStaffCountQuota()
	t.ReportSharingStatus = d.GetReportSharingStatus()
	t.ConstitutionStatus = d.GetConstitutionStatus()
	t.Overdue = d.GetOverdue()
	t.Rev = d.GetRev()
	t.SafePhone = d.GetSafePhone()
	t.CreatedAt = d.GetCreatedAt()
}

// db->domain
func (t *TenantEntity) ToDomainTenantEntity() domain.TenantEntityIntf {
	if t == nil {
		return nil
	}

	p := domain.TenantEntity{
		TenantID:            t.TenantID,
		LogoUrl:             t.LogoUrl,
		StoreName:           t.StoreName,
		Province:            t.Province,
		City:                t.City,
		District:            t.District,
		Street:              t.Street,
		ContactName:         t.ContactName,
		ContactPhone:        t.ContactPhone,
		SocialCreditCode:    t.SocialCreditCode,
		BusinessLicenseUrl:  t.BusinessLicenseUrl,
		StaffCountQuota:     t.StaffCountQuota,
		ReportSharingStatus: t.ReportSharingStatus,
		ConstitutionStatus:  t.ConstitutionStatus,
		Overdue:             t.Overdue,
		Rev:                 t.Rev,
		SafePhone:           t.SafePhone,
		CreatedAt:           t.CreatedAt,
	}
	return &p
}
