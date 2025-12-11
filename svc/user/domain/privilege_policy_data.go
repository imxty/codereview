package domain

func (p *PrivilegePolicyData) GetDataPrivilegeID() string {
	if p == nil {
		return ""
	}
	return p.DataPrivilegeID
}

func (p *PrivilegePolicyData) GetPrivilegeID() string {
	if p == nil {
		return ""
	}
	return p.PrivilegeID
}

func (p *PrivilegePolicyData) GetOrganizationID() string {
	if p == nil {
		return ""
	}
	return p.OrganizationID
}

func (p *PrivilegePolicyData) GetRev() int32 {
	if p == nil {
		return 0
	}
	return p.Rev
}

type PrivilegePolicyData struct {
	DataPrivilegeID string
	PrivilegeID     string
	OrganizationID  string
	Rev             int32
}

type PrivilegePolicyDataMapper interface {
	ToDomainPrivilegePolicyDataMapper
	FromDomainPrivilegePolicyDataMapper
}

type ToDomainPrivilegePolicyDataMapper interface {
	ToDomainPrivilegePolicyData() PrivilegePolicyDataIntf
}

type FromDomainPrivilegePolicyDataMapper interface {
	FromDomainPrivilegePolicyData(PrivilegePolicyDataIntf)
}

type PrivilegePolicyDataIntf interface {
	GetDataPrivilegeID() string
	GetPrivilegeID() string
	GetOrganizationID() string
	GetRev() int32
}

var _ PrivilegePolicyDataIntf = (*PrivilegePolicyData)(nil)
