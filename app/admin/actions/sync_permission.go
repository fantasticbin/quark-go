package actions

import (
	"strings"

	"github.com/gobeam/stringy"
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/model"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/message"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/actions"
	"gorm.io/gorm"
)

type SyncPermissionAction struct {
	actions.Action
}

// 同步权限，SyncPermission() | SyncPermission("同步权限")
func SyncPermission(options ...interface{}) *SyncPermissionAction {
	action := &SyncPermissionAction{}

	// 文字
	action.Name = "同步权限"
	if len(options) == 1 {
		action.Name = options[0].(string)
	}

	return action
}

// 初始化
func (p *SyncPermissionAction) Init(ctx *quark.Context) interface{} {

	// 执行成功后刷新的组件
	p.Reload = "table"

	// 是否具有loading，当action 的作用类型为ajax,submit时有效
	p.WithLoading = true

	// 设置展示位置
	p.SetOnlyOnIndex(true)

	// 行为类型
	p.ActionType = "ajax"

	return p
}

// 执行行为句柄
func (p *SyncPermissionAction) Handle(ctx *quark.Context, query *gorm.DB) error {
	permissions := ctx.Engine.GetUrlPaths()
	data := []model.Permission{}

	var names []string
	db.Client.Model(&model.Permission{}).Pluck("name", &names)
	for _, v := range permissions {
		if strings.Contains(v.Url, "/api/admin") {
			has := false
			hasPermission := false
			url := strings.ReplaceAll(v.Url, "/api/admin/", "")
			url = strings.ReplaceAll(url, "/", "_") + "_" + strings.ToLower(v.Method)
			name := stringy.
				New(url).
				CamelCase("?", "")

			// 判断数据库中是否已存在
			for _, nv := range names {
				if nv == name {
					has = true
				}
			}

			// 判断当前同步中是否已存在
			for _, pv := range data {
				if pv.Name == name {
					hasPermission = true
				}
			}

			if !has && !hasPermission {
				permission := model.Permission{
					Name:      name,
					Method:    v.Method,
					Path:      v.Url,
					GuardName: "admin",
				}
				data = append(data, permission)
			}
		}
	}
	if len(data) == 0 {
		return ctx.JSON(200, message.Error("暂无新增权限！"))
	}

	err := query.Create(data).Error
	if err != nil {
		return ctx.JSON(200, message.Error(err.Error()))
	}

	return ctx.JSON(200, message.Success("操作成功"))
}
