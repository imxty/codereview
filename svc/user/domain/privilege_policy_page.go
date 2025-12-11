package domain

func (p *PrivilegePolicyPage) GetPrivilegeID() string {
	if p == nil {
		return ""
	}
	return p.PrivilegeID
}

func (p *PrivilegePolicyPage) GetPresetID() string {
	if p == nil {
		return ""
	}
	return p.PresetID
}

type PrivilegePolicyPage struct {
	PrivilegeID string
	PresetID    string
}

type PrivilegePolicyPageMapper interface {
	ToDomainPrivilegePolicyPageMapper
	FromDomainPrivilegePolicyPageMapper
}

type ToDomainPrivilegePolicyPageMapper interface {
	ToDomainPrivilegePolicyPage() PrivilegePolicyPageIntf
}

type FromDomainPrivilegePolicyPageMapper interface {
	FromDomainPrivilegePolicyPage(PrivilegePolicyPageIntf)
}

type PrivilegePolicyPageIntf interface {
	GetPrivilegeID() string
	GetPresetID() string
}

var _ PrivilegePolicyPageIntf = (*PrivilegePolicyPage)(nil)
