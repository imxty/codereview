package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

// User 用户信息
type User struct {
	UserID                   string    `gorm:"primary_key;column:user_id"`
	OrganizationID           string    `gorm:"column:organization_id"`
	OrganizationContactPhone string    `gorm:"column:organization_contact_phone"`
	TenantID                 string    `gorm:"column:tenant_id"`
	Username                 string    `gorm:"column:username"`
	RoleType                 int32     `gorm:"column:role_type"`
	Nickname                 string    `gorm:"column:nickname"`
	Phone                    string    `gorm:"column:phone"`
	HashedPassword           string    `gorm:"column:hashed_password"`
	TenantLimit              int32     `gorm:"column:tenant_limit"`
	PrivilegeID              string    `gorm:"column:privilege_id"`
	Rev                      int32     `gorm:"column:rev"`
	Remark                   string    `gorm:"column:remark"`
	IsActivated              bool      `gorm:"column:is_activated"`
	CreatedAt                time.Time // 创建时间
	UpdatedAt                time.Time // 更新时间
	DeletedAt                *gorm.DeletedAt
}

func (u User) TableName() string {
	return "user"
}

// domain->db
func (u *User) FromDomainUser(d domain.UserIntf) {
	if u == nil || d == nil {
		return
	}

	u.UserID = d.GetUserID()
	u.OrganizationID = d.GetOrganizationID()
	u.OrganizationContactPhone = d.GetOrganizationContactPhone()
	u.TenantID = d.GetTenantID()
	u.Username = d.GetUsername()
	u.RoleType = d.GetRoleType()
	u.Nickname = d.GetNickname()
	u.Phone = d.GetPhone()
	u.HashedPassword = d.GetHashedPassword()
	u.TenantLimit = d.GetTenantLimit()
	u.PrivilegeID = d.GetPrivilegeID()
	u.Rev = d.GetRev()
	u.Remark = d.GetRemark()
	u.IsActivated = d.GetIsActivated()
}

// db->domain
func (u *User) ToDomainUser() domain.UserIntf {
	if u == nil {
		return nil
	}

	p := domain.User{
		UserID:                   u.UserID,
		OrganizationID:           u.OrganizationID,
		OrganizationContactPhone: u.OrganizationContactPhone,
		TenantID:                 u.TenantID,
		Username:                 u.Username,
		RoleType:                 u.RoleType,
		Nickname:                 u.Nickname,
		Phone:                    u.Phone,
		HashedPassword:           u.HashedPassword,
		TenantLimit:              u.TenantLimit,
		PrivilegeID:              u.PrivilegeID,
		Rev:                      u.Rev,
		Remark:                   u.Remark,
		IsActivated:              u.IsActivated,
	}
	return &p
}
