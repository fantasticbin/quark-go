package actions

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/action"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/actions"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
)

type EditModalAction struct {
	actions.Modal
}

// 编辑-弹窗类型，EditModal() | EditModal("编辑")
func EditModal(options ...interface{}) *EditModalAction {
	action := &EditModalAction{}

	// 文字
	action.Name = "编辑"
	if len(options) == 1 {
		action.Name = options[0].(string)
	}

	return action
}

// 初始化
func (p *EditModalAction) Init(ctx *quark.Context) interface{} {

	// 类型
	p.Type = "link"

	// 设置按钮大小,large | middle | small | default
	p.Size = "small"

	// 关闭时销毁 Modal 里的子元素
	p.DestroyOnClose = true

	// 执行成功后刷新的组件
	p.Reload = "table"

	// 设置展示位置
	p.SetOnlyOnIndexTableRow(true)

	return p
}

// 内容
func (p *EditModalAction) GetBody(ctx *quark.Context) interface{} {
	template := ctx.Template.(types.Resourcer)

	// 更新表单的接口
	api := template.UpdateApi(ctx)

	// 编辑页面获取表单数据接口
	initApi := template.EditValueApi(ctx)

	// 包裹在组件内的编辑页字段
	fields := template.UpdateFieldsWithinComponents(ctx)

	// 返回数据
	return (&form.Component{}).
		Init().
		SetStyle(map[string]interface{}{
			"paddingTop": "24px",
		}).
		SetKey("editModalForm", false).
		SetApi(api).
		SetInitApi(initApi).
		SetBody(fields).
		SetLabelCol(map[string]interface{}{
			"span": 6,
		}).
		SetWrapperCol(map[string]interface{}{
			"span": 18,
		})
}

// 弹窗行为
func (p *EditModalAction) GetActions(ctx *quark.Context) []interface{} {

	return []interface{}{
		(&action.Component{}).
			Init().
			SetLabel("取消").
			SetActionType("cancel"),

		(&action.Component{}).
			Init().
			SetLabel("提交").
			SetWithLoading(true).
			SetReload("table").
			SetActionType("submit").
			SetType("primary", false).
			SetSubmitForm("editModalForm"),
	}
}
