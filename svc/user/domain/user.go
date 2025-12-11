package domain

const (
	// 租户员工
	RoleTypeStaff = 0
	// 组织
	RoleTypeOrganization = 1
	// 商户
	RoleTypeTenant = 2
	// 运营后台管理员
	RoleTypeJmAdmin = 3
	// 管理员
	RoleTypeBoss = 4
)

type User struct {
	UserID                   string
	OrganizationID           string
	OrganizationContactPhone string
	TenantID                 string
	Username                 string
	RoleType                 int32
	Nickname                 string
	Phone                    string
	HashedPassword           string
	TenantLimit              int32
	PrivilegeID              string
	Rev                      int32
	Remark                   string
	IsActivated              bool
}

type UserMapper interface {
	ToDomainUserMapper
	FromDomainUserMapper
}

type ToDomainUserMapper interface {
	ToDomainUser() UserIntf
}

type FromDomainUserMapper interface {
	FromDomainUser(UserIntf)
}

type UserIntf interface {
	GetUserID() string
	GetOrganizationID() string
	GetOrganizationContactPhone() string
	GetTenantID() string
	GetUsername() string
	GetRoleType() int32
	GetNickname() string
	GetPhone() string
	GetHashedPassword() string
	GetTenantLimit() int32
	GetPrivilegeID() string
	GetRev() int32
	GetRemark() string
	GetIsActivated() bool
}

var _ UserIntf = (*User)(nil)

func (u *User) GetRemark() string {
	if u == nil {
		return ""
	}
	return u.Remark
}

func (u *User) GetUserID() string {
	if u == nil {
		return ""
	}
	return u.UserID
}

func (u *User) GetOrganizationID() string {
	if u == nil {
		return ""
	}
	return u.OrganizationID
}

func (u *User) GetOrganizationContactPhone() string {
	if u == nil {
		return ""
	}
	return u.OrganizationContactPhone
}

func (u *User) GetTenantID() string {
	if u == nil {
		return ""
	}
	return u.TenantID
}

func (u *User) GetUsername() string {
	if u == nil {
		return ""
	}
	return u.Username
}

func (u *User) GetRoleType() int32 {
	if u == nil {
		return 0
	}
	return u.RoleType
}

func (u *User) GetNickname() string {
	if u == nil {
		return ""
	}
	return u.Nickname
}

func (u *User) GetPhone() string {
	if u == nil {
		return ""
	}
	return u.Phone
}

func (u *User) GetHashedPassword() string {
	if u == nil {
		return ""
	}
	return u.HashedPassword
}

func (u *User) GetTenantLimit() int32 {
	if u == nil {
		return 0
	}
	return u.TenantLimit
}

func (u *User) GetPrivilegeID() string {
	if u == nil {
		return ""
	}
	return u.PrivilegeID
}

func (u *User) GetRev() int32 {
	if u == nil {
		return 0
	}
	return u.Rev
}

func (u *User) GetIsActivated() bool {
	if u == nil {
		return false
	}
	return u.IsActivated
}
