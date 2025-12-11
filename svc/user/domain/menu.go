package domain

import "github.com/rs/xid"

func (m *Menu) GetMenuID() xid.ID {
	if m == nil {
		return xid.NilID()
	}
	return m.MenuID
}

func (m *Menu) GetTitle() string {
	if m == nil {
		return ""
	}
	return m.Title
}

func (m *Menu) GetPath() string {
	if m == nil {
		return ""
	}
	return m.Path
}

func (m *Menu) GetSort() int32 {
	if m == nil {
		return 0
	}
	return m.Sort
}

func (m *Menu) GetVisible() bool {
	if m == nil {
		return false
	}
	return m.Visible
}

type Menu struct {
	MenuID  xid.ID
	Title   string
	Path    string
	Sort    int32
	Visible bool
}

type MenuMapper interface {
	ToDomainMenuMapper
	FromDomainMenuMapper
}

type ToDomainMenuMapper interface {
	ToDomainMenu() MenuIntf
}

type FromDomainMenuMapper interface {
	FromDomainMenu(MenuIntf)
}

type MenuIntf interface {
	GetMenuID() xid.ID
	GetTitle() string
	GetPath() string
	GetSort() int32
	GetVisible() bool
}

var _ MenuIntf = (*Menu)(nil)
