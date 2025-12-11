package casbin_test

import (
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/permission/casbin"
)

type testModel struct {
	name        string
	privilegeId string
	action      string
	resource    string
	want        bool
}

func TestMatch(t *testing.T) {
	p := []*casbin.PermissionPreset{
		{
			// 策略集 ID
			PresetId: "1",
			// 动作
			Actions: []string{"com.test.com/DeviceAPI.GetDevice"},
			// * 客体
			// action 作用的对象
			Resources: []string{"*"},
			// Allow / Deny
			Effect: "deny",
		},
		{
			// 策略集 ID
			PresetId: "1",
			// 动作
			Actions: []string{"com.test.com/DeviceAPI.*"},
			// * 客体
			// action 作用的对象
			Resources: []string{"*"},
			// Allow / Deny
			Effect: "allow",
		},
		{
			// 策略集 ID
			// 动作
			Actions: []string{"com.test.com/ReportAPI.GetReport", "com.test.com/ReportAPI.SubmitReport"},
			// * 客体
			// action 作用的对象
			Resources: []string{"*"},
			// Allow / Deny
			Effect: "allow",
		},
		{
			// 策略集 ID
			PresetId: "3",
			// 动作
			Actions: []string{"com.test.com/UserAPI*"},
			// * 客体
			// action 作用的对象
			Resources: []string{"*"},
			// Allow / Deny
			Effect: "allow",
		},
		{
			// 策略集 ID
			PresetId: "3",
			// 动作
			Actions: []string{"com.test.com/UserAPI.GetUser", "com.test.com/UserAPI.DeleteUser"},
			// * 客体
			// action 作用的对象
			Resources: []string{"*"},
			// Allow / Deny
			Effect: "deny",
		},
	}
	s := []*casbin.PolicyGroup{
		{
			// privilege_id
			PrivilegeId: "jinmu",
			// 策略集 ID
			PresetId: "1",
		},
		{
			PrivilegeId: "admin",
			PresetId:    "1",
		},
		{
			PrivilegeId: "admin",
			PresetId:    "2",
		},
		{
			PrivilegeId: "admin",
			PresetId:    "3",
		},
	}
	e, _ := casbin.NewCasbinPermission(p, s)

	tests := []testModel{
		{
			name:        "test1",
			privilegeId: "jinmu",
			resource:    "*",
			action:      "com.test.com/DeviceAPI.DeleteDevice",
			want:        true,
		},
		{
			name:        "test1",
			privilegeId: "jinmu",
			resource:    "*",
			action:      "com.test.com/DeviceAPI.GetDevice",
			want:        false,
		},
		{
			name:        "test1",
			privilegeId: "jinmu",
			resource:    "asdasdasd",
			action:      "com.test.com/DeviceAPI.GetDevice",
			want:        false,
		},
		{
			name:        "test1",
			privilegeId: "jinmu",
			resource:    "asdasdasd",
			action:      "com.test.com/DeviceAPI.DeleteDevice",
			want:        true,
		},
		{
			name:        "test1",
			privilegeId: "jinmu",
			resource:    "",
			action:      "com.test.com/DeviceAPI.DeleteDevice",
			want:        true,
		},
		{
			name:        "test1",
			privilegeId: "jinmu",
			resource:    "",
			action:      "com.test.com/DeviceAPI.GetDevice",
			want:        false,
		},
		{
			name:        "test1",
			privilegeId: "admin",
			resource:    "*",
			action:      "com.test.com/UserAPI.GetUser",
			want:        false,
		},
		{
			name:        "test1",
			privilegeId: "admin",
			resource:    "*",
			action:      "com.test.com/UserAPI.SignIn",
			want:        true,
		},
	}

	for _, v := range tests {
		// 验证
		t.Run(v.name, func(t *testing.T) {
			if ok, _ := e.CheckPermission(v.privilegeId, v.action, v.resource); ok != v.want {
				t.Errorf("check permission failed, %s %s %s", v.privilegeId, v.resource, v.action)
			}
		})
	}
}

func TestRemovePermission(t *testing.T) {
	p := []*casbin.PermissionPreset{
		{
			// 策略集 ID
			PresetId: "1",
			// 动作
			Actions: []string{"com.test.com/DeviceAPI.GetDevice"},
			// * 客体
			// action 作用的对象
			Resources: []string{"*"},
			// Allow / Deny
			Effect: "deny",
		},
		{
			PresetId:  "1",
			Actions:   []string{"com.test.com/DeviceAPI.*"},
			Resources: []string{"*"},
			Effect:    "allow",
		},
		{
			PresetId:  "2",
			Actions:   []string{"com.test.com/ReportAPI.GetReport", "com.test.com/ReportAPI.SubmitReport"},
			Resources: []string{"*"},
			Effect:    "allow",
		},
		{
			PresetId:  "3",
			Actions:   []string{"com.test.com/UserAPI*"},
			Resources: []string{"*"},
			Effect:    "allow",
		},
		{
			PresetId:  "3",
			Actions:   []string{"com.test.com/UserAPI.GetUser", "com.test.com/UserAPI.DeleteUser"},
			Resources: []string{"*"},
			Effect:    "deny",
		},
	}
	s := []*casbin.PolicyGroup{
		{
			// privilege_id
			PrivilegeId: "jinmu",
			// 策略集 ID
			PresetId: "1",
		},
		{
			PrivilegeId: "admin",
			PresetId:    "1",
		},
		{
			PrivilegeId: "admin",
			PresetId:    "2",
		},
		{
			PrivilegeId: "admin",
			PresetId:    "3",
		},
	}

	e, _ := casbin.NewCasbinPermission(p, s)

	e.RemovePermission("1")

	tests := []testModel{
		{
			name:        "test1",
			privilegeId: "jinmu",
			resource:    "*",
			action:      "com.test.com/DeviceAPI.DeleteDevice",
			want:        false,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if ok, _ := e.CheckPermission(v.privilegeId, v.action, v.resource); ok != v.want {
				t.Errorf("check permission failed, %s %s %s", v.privilegeId, v.resource, v.action)
			}
		})
	}
}

func TestRemoveAppPermissionGroup(t *testing.T) {
	p := []*casbin.PermissionPreset{
		{
			// 策略集ID
			PresetId: "1",
			// 动作
			Actions: []string{"com.test.com/DeviceAPI.GetDevice"},
			// * 客体
			// action作用的对象
			Resources: []string{"*"},
			// Allow/Deny
			Effect: "deny",
		},
		{
			// 策略集ID
			PresetId: "1",
			// 动作
			Actions: []string{"com.test.com/DeviceAPI.*"},
			// * 客体
			// action作用的对象
			Resources: []string{"*"},
			// Allow/Deny
			Effect: "allow",
		},
		{
			// 策略集ID
			PresetId: "2",
			// 动作
			Actions: []string{"com.test.com/ReportAPI.GetReport", "com.test.com/ReportAPI.SubmitReport"},
			// * 客体
			// action作用的对象
			Resources: []string{"*"},
			// Allow/Deny
			Effect: "allow",
		},
		{
			// 策略集ID
			PresetId: "3",
			// 动作
			Actions: []string{"com.test.com/UserAPI*"},
			// * 客体
			// action作用的对象
			Resources: []string{"*"},
			// Allow/Deny
			Effect: "allow",
		},
		{
			// 策略集ID
			PresetId: "3",
			// 动作
			Actions: []string{"com.test.com/UserAPI.GetUser", "com.test.com/UserAPI.DeleteUser"},
			// * 客体
			// action作用的对象
			Resources: []string{"*"},
			// Allow/Deny
			Effect: "deny",
		},
	}
	s := []*casbin.PolicyGroup{
		{
			// privilege_id
			PrivilegeId: "jinmu",
			// 策略集ID
			PresetId: "1",
		},
		{
			// app_id
			PrivilegeId: "admin",
			// 策略集ID
			PresetId: "1",
		},
		{
			// app_id
			PrivilegeId: "admin",
			// 策略集ID
			PresetId: "2",
		},
		{
			// app_id
			PrivilegeId: "admin",
			// 策略集ID
			PresetId: "3",
		},
	}

	e, _ := casbin.NewCasbinPermission(p, s)

	e.DeletePrivilegePermissionGroup("admin")

	tests := []testModel{
		{
			name:        "test1",
			privilegeId: "admin",
			resource:    "*",
			action:      "com.test.com/DeviceAPI.DeleteDevice",
			want:        false,
		},
	}

	for _, v := range tests {
		// 验证
		t.Run(v.name, func(t *testing.T) {
			if ok, _ := e.CheckPermission(v.privilegeId, v.action, v.resource); ok != v.want {
				t.Errorf("check permission failed,%s %s %s", v.privilegeId, v.resource, v.action)
			}
		})
	}
}
