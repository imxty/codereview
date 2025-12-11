package casbin

import (
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/jinmukeji/huimaibao-service/pkg/permission"
	"github.com/jinmukeji/plat-pkg/v4/experiments/whitelist/matcher"
)

var (
	// com.shangyikangyou.huimaibao.service/ReportAPI.*
	// 分解为 namespace: com.shangyikangyou.huimaibao.service
	// API: ReportAPI.*
	reg = regexp.MustCompile(`(\S*)\/(\S*)`)
)

// 基于 casbin 实现的 Permission 管理
type CasbinPermission struct {
	// casbin enforcer
	*casbin.Enforcer
	// 读写锁
	mutex sync.RWMutex
}

var _ permission.Permission = (*CasbinPermission)(nil)

// 初始化 CasbinPermission
func NewCasbinPermission(ps []*PermissionPreset, aps []*PolicyGroup) (*CasbinPermission, error) {
	// 读取 model
	// rbac_with_deny_model
	// https://github.com/casbin/casbin/blob/master/examples/rbac_with_deny_model.conf
	m, err := model.NewModelFromString(`
        [request_definition]
        r = sub, obj, act
        
        [policy_definition]
        p = sub, obj, act, eft
        
        [role_definition]
        g = _,_
        
        [matchers]
        m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && match_act(r.act, p.act)
        
        [policy_effect]
        e = some(where (p.eft == allow)) && !some(where (p.eft == deny)) 
    `)
	if err != nil {
		return nil, err
	}

	// 初始化 enforcer
	e, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, err
	}

	// 添加自定义方法
	e.AddFunction("match_act", keyMatchFunc)

	c := &CasbinPermission{
		Enforcer: e,
	}

	// 构建 casbin policy
	for _, v := range ps {
		// 添加 Permission
		err := c.AddPermissions(v.PresetId, v.Effect, v.Actions, v.Resources)
		if err != nil {
			return nil, err
		}
	}

	// 构建 casbin Group
	for i := 0; i < len(aps); i++ {
		_, err := e.AddGroupingPolicy(aps[i].PrivilegeId, aps[i].PresetId)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// CheckPermission 检测权限
func (c *CasbinPermission) CheckPermission(privilegeId, action, resource string) (bool, error) {
	// 加上读锁
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Enforce(privilegeId, resource, action)
}

// RemovePermission 删除权限
func (c *CasbinPermission) RemovePermission(presetId string) error {
	// 加上写锁
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 删除权限
	_, err := c.RemoveFilteredPolicy(0, presetId)
	if err != nil {
		return fmt.Errorf("remove policy failed: %w", err)
	}
	// 把所有与该 preset 有关的权限组也删除
	_, err = c.RemoveFilteredGroupingPolicy(1, presetId)
	if err != nil {
		return fmt.Errorf("remove app group preset failed: %w", err)
	}
	return nil
}

// AddPermissions 添加权限
func (c *CasbinPermission) AddPermissions(presetId, effect string, actions, resources []string) error {
	// 加上写锁
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 构建 policies
	policies := make([][]string, len(actions)*len(resources))
	index := 0
	for _, v := range resources {
		for _, j := range actions {
			policies[index] = []string{presetId, v, j, effect}
			index++
		}
	}
	_, err := c.AddPolicies(policies)
	if err != nil {
		return fmt.Errorf("add policy failed: %w", err)
	}
	return nil
}

// UpdatePermission 更新权限
func (c *CasbinPermission) UpdatePermission(presetId, effect string, actions, resources []string) error {
	// 加上写锁
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 先删除
	err := c.RemovePermission(presetId)
	if err != nil {
		return fmt.Errorf("update policy failed: %w", err)
	}
	// 再添加
	return c.AddPermissions(presetId, effect, actions, resources)
}

// Group 有关
// AddPermissionGroup 添加组关系
func (c *CasbinPermission) AddPermissionGroup(privilegeId, presetId string) error {
	// 加上写锁
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 添加组关系
	_, err := c.AddGroupingPolicy(privilegeId, presetId)
	if err != nil {
		return fmt.Errorf("add group policy failed: %w", err)
	}
	return nil
}

// DeletePermissionGroup 删除权限组
func (c *CasbinPermission) DeletePermissionGroup(privilegeId, presetId string) error {
	// 加上写锁
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 删除权限组
	_, err := c.RemoveGroupingPolicy(privilegeId, presetId)
	if err != nil {
		return fmt.Errorf("remove group policy failed: %w", err)
	}
	return nil
}

// DeletePrivilegePermissonGroup 删除特定一组权限组
func (c *CasbinPermission) DeletePrivilegePermissionGroup(privilegeId string) error {
	// 加上写锁
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 删除所有 privilegeId 的组
	_, err := c.RemoveFilteredGroupingPolicy(0, privilegeId)
	if err != nil {
		return fmt.Errorf("remove app group policy failed: %w", err)
	}
	return nil
}

// 用于 key 匹配的函数封装，内部为自定义实现的 key 匹配规则
func keyMatchFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("invalid args")
	}
	// key1 为 request 请求的 key
	// key2 为数据库保存的 key
	key1, ok := args[0].(string)
	if !ok {
		return nil, errors.New("failed to convert argument to string")
	}
	key2, ok := args[1].(string)
	if !ok {
		return nil, errors.New("failed to convert argument to string")
	}

	return (bool)(matchAct(key1, key2)), nil
}

// 作用于 casbin 的 model 定义的 matchers 中
// model 中定义 match_act(r.act, p.act)
// match_act 将会调用 keyMatchFunc 再由 keyMatchFunc 调用 matchAct 进行 key 的判断
func matchAct(act1, act2 string) bool {
	// key1 为 request 请求的 key
	// key2 为数据库保存的 key
	// sample key1 = com.shangyikangyou.huimaibao.service/ReportAPI.GetReport
	// key2 = com.shangyikangyou.huimaibao.service/ReportAPI.*
	// 解析 key2
	// com.shangyikangyou.huimaibao.service/ReportAPI.*
	// 分解为 namespace: com.shangyikangyou.huimaibao.service
	// API: ReportAPI.*
	res := reg.FindStringSubmatch(act1)
	if len(res) != 3 {
		return false
	}
	// 判断规则详见 plat-pkg/matcher
	return matcher.Match(act2, res[1], res[2])
}
