package domain

func (p *PrivilegePolicy) GetPrivilegeID() string {
	if p == nil {
		return ""
	}
	return p.PrivilegeID
}

func (p *PrivilegePolicy) GetPrivilegeName() string {
	if p == nil {
		return ""
	}
	return p.PrivilegeName
}

func (p *PrivilegePolicy) GetRemark() string {
	if p == nil {
		return ""
	}
	return p.Remark
}

func (p *PrivilegePolicy) GetRev() int32 {
	if p == nil {
		return 0
	}
	return p.Rev
}

type PrivilegePolicy struct {
	PrivilegeID   string
	PrivilegeName string
	Remark        string
	Rev           int32
}

type PrivilegePolicyMapper interface {
	ToDomainPrivilegePolicyMapper
	FromDomainPrivilegePolicyMapper
}

type ToDomainPrivilegePolicyMapper interface {
	ToDomainPrivilegePolicy() PrivilegePolicyIntf
}

type FromDomainPrivilegePolicyMapper interface {
	FromDomainPrivilegePolicy(PrivilegePolicyIntf)
}

type PrivilegePolicyIntf interface {
	GetPrivilegeID() string
	GetPrivilegeName() string
	GetRemark() string
	GetRev() int32
}

var _ PrivilegePolicyIntf = (*PrivilegePolicy)(nil)
