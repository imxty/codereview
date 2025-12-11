package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

// Menu 菜单
type Menu struct {
	MenuID  xid.ID `gorm:"primary_key;column:menu_id"`
	Title   string `gorm:"column:title"`
	Path    string `gorm:"column:path"`
	Sort    int32  `gorm:"column:sort"`
	Visible bool   `gorm:"column:visible"`
}

// PrivilegeMenu
type PrivilegeMenu struct {
	PrivilegeID string          `gorm:"column:privilege_id"`
	MenuID      string          `gorm:"column:menu_id"`
	Rev         int32           `gorm:"column:rev"`
	CreatedAt   time.Time       // 创建时间
	UpdatedAt   time.Time       // 更新时间
	DeletedAt   *gorm.DeletedAt // 删除时间
}

func (m Menu) TableName() string {
	return "menu"
}

func (m PrivilegeMenu) TableName() string {
	return "privilege_menu"
}

// domain->db
func (m *Menu) FromDomainMenu(d domain.MenuIntf) {
	if m == nil || d == nil {
		return
	}

	m.MenuID = d.GetMenuID()
	m.Title = d.GetTitle()
	m.Path = d.GetPath()
	m.Sort = d.GetSort()
	m.Visible = d.GetVisible()
}

// db->domain
func (m *Menu) ToDomainMenu() domain.MenuIntf {
	if m == nil {
		return nil
	}

	p := domain.Menu{
		MenuID:  m.MenuID,
		Title:   m.Title,
		Path:    m.Path,
		Sort:    m.Sort,
		Visible: m.Visible,
	}
	return &p
}
