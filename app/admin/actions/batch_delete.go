package actions

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/message"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/actions"
	"gorm.io/gorm"
)

type BatchDeleteAction struct {
	actions.Action
}

// 批量删除，BatchDelete() | BatchDelete("批量删除")
func BatchDelete(options ...interface{}) *BatchDeleteAction {
	action := &BatchDeleteAction{}

	action.Name = "批量删除"
	if len(options) == 1 {
		action.Name = options[0].(string)
	}

	return action
}

// 初始化
func (p *BatchDeleteAction) Init(ctx *quark.Context) interface{} {

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
func (p *BatchDeleteAction) GetApiParams() []string {
	return []string{
		"id",
	}
}

// 执行行为句柄
func (p *BatchDeleteAction) Handle(ctx *quark.Context, query *gorm.DB) error {
	err := query.Delete("").Error
	if err != nil {
		return ctx.JSON(200, message.Error(err.Error()))
	}

	return ctx.JSON(200, message.Success("操作成功"))
}
