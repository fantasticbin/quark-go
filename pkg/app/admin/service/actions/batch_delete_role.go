package actions

import (
	"strconv"
	"strings"

	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/component/message"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/model"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/template/resource/actions"
	"github.com/quarkcloudio/quark-go/v3/pkg/builder"
	"gorm.io/gorm"
)

type BatchDeleteRoleAction struct {
	actions.Action
}

// 批量删除角色
func BatchDeleteRole() *BatchDeleteRoleAction {
	return &BatchDeleteRoleAction{}
}

// 初始化
func (p *BatchDeleteRoleAction) Init(ctx *builder.Context) interface{} {

	// 设置按钮文字
	p.Name = "批量删除"

	// 设置按钮类型,primary | ghost | dashed | link | text | default
	p.Type = "link"

	// 设置按钮大小,large | middle | small | default
	p.Size = "small"

	//  执行成功后刷新的组件
	p.Reload = "table"

	// 当行为在表格行展示时，支持js表达式
	p.WithConfirm("确定要删除吗？", "删除后数据将无法恢复，请谨慎操作！", "modal")

	// 在表格多选弹出层展示
	p.SetOnlyOnIndexTableAlert(true)

	return p
}

// 行为接口接收的参数，当行为在表格行展示的时候，可以配置当前行的任意字段
func (p *BatchDeleteRoleAction) GetApiParams() []string {
	return []string{
		"id",
	}
}

// 执行行为句柄
func (p *BatchDeleteRoleAction) Handle(ctx *builder.Context, query *gorm.DB) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.JSON(200, message.Error("参数错误！"))
	}

	err := query.Delete("").Error
	if err != nil {
		return ctx.JSON(200, message.Error(err.Error()))
	}

	ids := strings.Split(id.(string), ",")
	if len(ids) > 0 {
		for _, v := range ids {
			idInt, err := strconv.Atoi(v)
			if err != nil {
				return ctx.JSON(200, message.Error(err.Error()))
			}

			// 清理casbin里的角色
			(&model.CasbinRule{}).RemoveRoleMenuAndPermissions(idInt)
		}
	} else {
		idInt, err := strconv.Atoi(id.(string))
		if err != nil {
			return ctx.JSON(200, message.Error(err.Error()))
		}

		// 清理casbin里的角色
		(&model.CasbinRule{}).RemoveRoleMenuAndPermissions(idInt)
	}

	return ctx.JSON(200, message.Success("操作成功"))
}
