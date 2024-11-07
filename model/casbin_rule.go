package model

import (
	"strconv"
	"strings"

	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/redis/go-redis/v9"
)

// 字段
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:100;uniqueIndex:unique_index"`
	V0    string `gorm:"size:100;uniqueIndex:unique_index"`
	V1    string `gorm:"size:100;uniqueIndex:unique_index"`
	V2    string `gorm:"size:100;uniqueIndex:unique_index"`
	V3    string `gorm:"size:100;uniqueIndex:unique_index"`
	V4    string `gorm:"size:100;uniqueIndex:unique_index"`
	V5    string `gorm:"size:100;uniqueIndex:unique_index"`
}

var Enforcer *casbin.Enforcer

// 获取Enforcer
func (p *CasbinRule) Enforcer() (enforcer *casbin.Enforcer, err error) {
	if Enforcer != nil {
		return Enforcer, err
	}

	a, err := gormadapter.NewAdapterByDBWithCustomTable(db.Client, &CasbinRule{}, "casbin_rules")
	if err != nil {
		return nil, err
	}
	m, err := casbinmodel.NewModelFromString(`
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`)
	if err != nil {
		return nil, err
	}

	Enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	redisConfig := quark.GetConfig().RedisConfig
	if redisConfig != nil {

		// 不同的Channel
		appKey := quark.GetConfig().AppKey
		w, _ := rediswatcher.NewWatcher(redisConfig.Host+":"+redisConfig.Port, rediswatcher.WatcherOptions{
			Options: redis.Options{
				Network:  "tcp",
				Password: redisConfig.Password,
			},
			Channel: "/casbin_" + appKey,
		})

		// Set the watcher for the enforcer.
		err = Enforcer.SetWatcher(w)
		if err != nil {
			return nil, err
		}

		// Or use the default callback
		err = w.SetUpdateCallback(func(s string) {
			Enforcer.LoadPolicy()
		})
		if err != nil {
			return nil, err
		}
	}

	return Enforcer, err
}

// 查看是否放行
func (p *CasbinRule) Enforce(sub string, obj string, act string) (result bool, err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return
	}

	result, err = enforcer.Enforce(sub, obj, act)

	return
}

// 添加菜单拥有的权限
func (p *CasbinRule) AddMenuPermission(menuId int, permissionIds interface{}) (err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return err
	}

	permissions, err := (&Permission{}).GetListByIds(permissionIds)
	if err != nil {
		return err
	}

	rules := [][]string{}
	addedRules := make(map[string]bool)

	// 菜单拥有的权限
	for _, v := range permissions {
		rule := "menu|" + strconv.Itoa(menuId) + v.Name + "MenuHasPermission"
		if !addedRules[rule] {
			rules = append(rules, []string{"menu|" + strconv.Itoa(menuId), v.Name, "MenuHasPermission"})
			addedRules[rule] = true
		}
	}

	p.RemoveMenuPermissions(menuId)
	if len(rules) == 0 {
		return
	}

	_, err = enforcer.AddPolicies(rules)
	if err != nil {
		return err
	}

	enforcer.SavePolicy()

	return
}

// 删除菜单拥有的权限
func (p *CasbinRule) RemoveMenuPermissions(menuId int) (err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return err
	}

	_, err = enforcer.DeleteUser("menu|" + strconv.Itoa(menuId))
	if err != nil {
		return err
	}

	enforcer.SavePolicy()

	return
}

// 获取菜单拥有的权限
func (p *CasbinRule) GetMenuPermissions(menuId int) (permissions []Permission, err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return
	}

	permissionNames := []string{}
	menuHasPermissions := enforcer.GetPermissionsForUser("menu|" + strconv.Itoa(menuId))
	for _, v := range menuHasPermissions {
		permissionNames = append(permissionNames, v[1])
	}

	permissions, err = (&Permission{}).GetListByNames(permissionNames)

	return
}

// 给角色添加菜单及权限
func (p *CasbinRule) AddMenuAndPermissionToRole(roleId int, menuIds []int) (err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return err
	}

	rules := [][]string{}
	addedRules := make(map[string]bool)

	// 角色拥有的菜单
	for _, v := range menuIds {
		rule := "role|" + strconv.Itoa(roleId) + "menu|" + strconv.Itoa(v) + "RoleHasMenu"
		if !addedRules[rule] {
			rules = append(rules, []string{"role|" + strconv.Itoa(roleId), "menu|" + strconv.Itoa(v), "RoleHasMenu"})
			addedRules[rule] = true
		}
	}

	// 角色拥有的权限
	for _, menuId := range menuIds {
		menuHasPermissions, err := p.GetMenuPermissions(menuId)
		if err == nil {
			for _, menuHasPermission := range menuHasPermissions {
				rule := "role|" + strconv.Itoa(roleId) + menuHasPermission.Path + menuHasPermission.Method
				if !addedRules[rule] {
					rules = append(rules, []string{"role|" + strconv.Itoa(roleId), menuHasPermission.Path, menuHasPermission.Method})
					addedRules[rule] = true
				}
			}
		}
	}

	// 先清理数据
	p.RemoveRoleMenuAndPermissions(roleId)

	if len(rules) > 0 {
		// 添加策略
		_, err = enforcer.AddPolicies(rules)
		if err != nil {
			return err
		}
	}

	return
}

// 删除角色拥有的菜单及权限
func (p *CasbinRule) RemoveRoleMenuAndPermissions(roleId int) (err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return err
	}

	_, err = enforcer.DeleteUser("role|" + strconv.Itoa(roleId))
	if err != nil {
		return err
	}

	return
}

// 获取角色拥有的菜单
func (p *CasbinRule) GetRoleMenus(roleId int) (menus []Menu, err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return
	}

	menuIds := []interface{}{}
	roleHasPermissions := enforcer.GetPermissionsForUser("role|" + strconv.Itoa(roleId))
	for _, v := range roleHasPermissions {
		if v[2] == "RoleHasMenu" {
			menuIdArr := strings.Split(v[1], "|")
			if len(menuIdArr) > 1 {
				menuIds = append(menuIds, menuIdArr[1])
			}
		}
	}
	menus, err = (&Menu{}).GetListByIds(menuIds)

	return
}

// 获取角色拥有的权限
func (p *CasbinRule) GetRolePermissions(roleId int) (permissions []Permission, err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return
	}

	permissionNames := []string{}
	roleHasPermissions := enforcer.GetPermissionsForUser("role|" + strconv.Itoa(roleId))
	for _, v := range roleHasPermissions {
		if v[2] != "RoleHasMenu" {
			permissionNames = append(permissionNames, v[1])
		}
	}
	permissions, err = (&Permission{}).GetListByNames(permissionNames)

	return
}

// 添加用户拥有的角色
func (p *CasbinRule) AddUserRole(modelId int, roleIds []int) (err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return err
	}

	roles := []string{}
	addedRules := make(map[string]bool)

	for _, v := range roleIds {
		rule := "role|" + strconv.Itoa(v)
		if !addedRules[rule] {
			roles = append(roles, "role|"+strconv.Itoa(v))
			addedRules[rule] = true
		}
	}

	p.RemoveUserRoles(modelId)

	if len(roles) > 0 {
		_, err = enforcer.AddRolesForUser("admin|"+strconv.Itoa(modelId), roles)
		if err != nil {
			return err
		}
	}

	return
}

// 删除用户拥有的角色
func (p *CasbinRule) RemoveUserRoles(modelId int) (err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return err
	}

	_, err = enforcer.DeleteRolesForUser("admin|" + strconv.Itoa(modelId))
	if err != nil {
		return err
	}

	return
}

// 获取用户拥有的角色
func (p *CasbinRule) GetUserRoles(modelId int) (roles []Role, err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return
	}

	roleStrIds, err := enforcer.GetRolesForUser("admin|" + strconv.Itoa(modelId))
	if err != nil {
		return
	}

	roleIds := []interface{}{}
	for _, v := range roleStrIds {
		roleIdArr := strings.Split(v, "|")
		if len(roleIdArr) > 1 {
			roleIds = append(roleIds, roleIdArr[1])
		}
	}
	roles, err = (&Role{}).GetListByIds(roleIds)

	return
}

// 获取用户拥有的菜单
func (p *CasbinRule) GetUserMenus(modelId int) (menus []Menu, err error) {
	getMenus := []Menu{}
	roles, err := p.GetUserRoles(modelId)
	if err != nil {
		return
	}

	for _, v := range roles {
		roleMenus, err := p.GetRoleMenus(v.Id)
		if err == nil {
			getMenus = append(getMenus, roleMenus...)
		}
	}

	return getMenus, nil
}

// 给角色添加数据权限
func (p *CasbinRule) AddDepartmentToRole(roleId int, departmentIds []int) (err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return err
	}

	rules := [][]string{}
	addedRules := make(map[string]bool)

	// 角色拥有的菜单
	for _, v := range departmentIds {
		rule := "roleDepartment|" + strconv.Itoa(roleId) + strconv.Itoa(v) + "RoleHasDepartment"
		if !addedRules[rule] {
			rules = append(rules, []string{"roleDepartment|" + strconv.Itoa(roleId), strconv.Itoa(v), "RoleHasDepartment"})
			addedRules[rule] = true
		}
	}

	// 先清理数据
	p.RemoveRoleDepartments(roleId)

	if len(rules) > 0 {
		// 添加策略
		_, err = enforcer.AddPolicies(rules)
		if err != nil {
			return err
		}
	}

	return
}

// 删除角色拥有的部门
func (p *CasbinRule) RemoveRoleDepartments(roleId int) (err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return err
	}

	_, err = enforcer.DeleteUser("roleDepartment|" + strconv.Itoa(roleId))
	if err != nil {
		return err
	}

	return
}

// 获取角色拥有的部门Ids
func (p *CasbinRule) GetRoleDepartmentIds(roleId int) (ids []int, err error) {
	enforcer, err := p.Enforcer()
	if err != nil {
		return
	}

	roleHasDepartmentIds := enforcer.GetPermissionsForUser("roleDepartment|" + strconv.Itoa(roleId))
	for _, v := range roleHasDepartmentIds {
		departmentId, err := strconv.Atoi(v[1])
		if err != nil {
			return nil, err
		}
		ids = append(ids, departmentId)
	}

	return
}

// 获取角色拥有的部门
func (p *CasbinRule) GetRoleDepartments(roleId int) (departments []Department, err error) {
	departmentIds, err := p.GetRoleDepartmentIds(roleId)
	if err != nil {
		return
	}
	return (&Department{}).GetListByIds(departmentIds)
}

// 获取用户拥有的部门
func (p *CasbinRule) GetUserDepartments(modelId int) (menus []Department, err error) {
	getDepartments := []Department{}
	roles, err := p.GetUserRoles(modelId)
	if err != nil {
		return
	}

	userInfo, err := (&User{}).GetInfoById(modelId)
	if err != nil {
		return
	}

	for _, v := range roles {
		switch v.DataScope {
		case 1:
			departments, err := (&Department{}).GetList()
			if err == nil {
				getDepartments = append(getDepartments, departments...)
			}
		case 2:
			departments, err := p.GetRoleDepartments(v.Id)
			if err == nil {
				getDepartments = append(getDepartments, departments...)
			}
		case 3:
			department, err := (&Department{}).GetInfoById(userInfo.DepartmentId)
			if err == nil {
				getDepartments = append(getDepartments, department)
			}
		case 4:
			department, err := (&Department{}).GetInfoById(userInfo.DepartmentId)
			if err == nil {
				getDepartments = append(getDepartments, department)
			}
			departments := (&Department{}).GetChildrenDepartments(userInfo.DepartmentId)
			getDepartments = append(getDepartments, departments...)
		case 5:
		}
	}

	return getDepartments, nil
}
