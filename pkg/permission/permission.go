package permission

type Permission interface {
	// Policy 有关
	// CheckPermission 检测权限
	CheckPermission(privilegeId, action, resource string) (bool, error)
	// RemovePermission 删除权限
	RemovePermission(persetId string) error
	// AddPermissions 添加权限
	AddPermissions(persetId, effect string, actions, resources []string) error
	// UpdatePermission 更新权限
	UpdatePermission(persetId, effect string, actions, resources []string) error

	// Group 有关
	// AddPermissionGroup 添加组关系
	AddPermissionGroup(privilegeId, presetId string) error
	// DeletePermissionGroup 删除权限组
	DeletePermissionGroup(privilegeId, presetId string) error
	// DeletePrivilegePermissionGroup 删除特定一组权限组
	DeletePrivilegePermissionGroup(privilegeId string) error
}
