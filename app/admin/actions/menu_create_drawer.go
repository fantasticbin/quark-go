package actions

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/action"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/actions"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
)

type MenuCreateDrawerAction struct {
	actions.Drawer
}

// 创建菜单-抽屉类型
func MenuCreateDrawer() *MenuCreateDrawerAction {
	return &MenuCreateDrawerAction{}
}

// 初始化
func (p *MenuCreateDrawerAction) Init(ctx *quark.Context) interface{} {
	template := ctx.Template.(types.Resourcer)

	// 文字
	p.Name = "创建" + template.GetTitle()

	// 类型
	p.Type = "primary"

	// 图标
	p.Icon = "plus-circle"

	// 执行成功后刷新的组件
	p.Reload = "table"

	// 关闭时销毁 Drawer 里的子元素
	p.DestroyOnClose = true

	// 抽屉弹出层宽度
	p.Width = 750

	// 设置展示位置
	p.SetOnlyOnIndex(true)

	return p
}

// 内容
func (p *MenuCreateDrawerAction) GetBody(ctx *quark.Context) interface{} {
	template := ctx.Template.(types.Resourcer)

	// 包裹在组件内的编辑页字段
	api := template.CreationApi(ctx)

	// 包裹在组件内的创建页字段
	fields := template.CreationFieldsWithinComponents(ctx)

	// 创建页面显示前回调
	data := template.BeforeCreating(ctx)

	// 返回数据
	return (&form.Component{}).
		Init().
		SetKey("createDrawerForm", false).
		SetLayout("vertical").
		SetApi(api).
		SetBody(fields).
		SetInitialValues(data)
}

// 弹窗行为
func (p *MenuCreateDrawerAction) GetActions(ctx *quark.Context) []interface{} {

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
			SetSubmitForm("createDrawerForm"),
	}
}
