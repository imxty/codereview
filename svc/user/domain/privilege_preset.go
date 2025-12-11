package domain

func (p *PrivilegePreset) GetPresetID() string {
	if p == nil {
		return ""
	}
	return p.PresetID
}

func (p *PrivilegePreset) GetPresetName() string {
	if p == nil {
		return ""
	}
	return p.PresetName
}

func (p *PrivilegePreset) GetPolicies() string {
	if p == nil {
		return ""
	}
	return p.Policies
}

func (p *PrivilegePreset) GetRemark() string {
	if p == nil {
		return ""
	}
	return p.Remark
}

func (p *PrivilegePreset) GetRev() int32 {
	if p == nil {
		return 0
	}
	return p.Rev
}

type PrivilegePreset struct {
	PresetID   string
	PresetName string
	Policies   string
	Remark     string
	Rev        int32
}

type PrivilegePresetMapper interface {
	ToDomainPrivilegePresetMapper
	FromDomainPrivilegePresetMapper
}

type ToDomainPrivilegePresetMapper interface {
	ToDomainPrivilegePreset() PrivilegePresetIntf
}

type FromDomainPrivilegePresetMapper interface {
	FromDomainPrivilegePreset(PrivilegePresetIntf)
}

type PrivilegePresetIntf interface {
	GetPresetID() string
	GetPresetName() string
	GetPolicies() string
	GetRemark() string
	GetRev() int32
}

var _ PrivilegePresetIntf = (*PrivilegePreset)(nil)
