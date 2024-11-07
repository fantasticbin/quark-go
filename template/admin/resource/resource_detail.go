package resource

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/card"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/tabs"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
)

// 详情页标题
func (p *Template) DetailTitle(ctx *quark.Context) string {
	template := ctx.Template.(types.Resourcer)
	title := template.GetTitle()

	return title + "详情"
}

// 渲染详情页组件
func (p *Template) DetailComponentRender(ctx *quark.Context, data map[string]interface{}) interface{} {

	// 详情页标题
	title := p.DetailTitle(ctx)

	// 详情页右上角自定义区域行为
	formExtraActions := p.DetailExtraActions(ctx)

	// 包裹在组件内的详情页字段
	fields := p.DetailFieldsWithinComponents(ctx, data)

	// 包裹在组件内的详情页字段
	formActions := p.DetailActions(ctx)

	return p.DetailWithinCard(
		ctx,
		title,
		formExtraActions,
		fields,
		formActions,
		data,
	)
}

// 在卡片内的详情页组件
func (p *Template) DetailWithinCard(
	ctx *quark.Context,
	title string,
	extra interface{},
	fields interface{},
	actions []interface{},
	data map[string]interface{}) interface{} {

	return (&card.Component{}).
		Init().
		SetTitle(title).
		SetHeaderBordered(true).
		SetExtra(extra).
		SetBody(fields)
}

// 在标签页内的详情页组件
func (p *Template) DetailWithinTabs(
	ctx *quark.Context,
	title string,
	extra interface{},
	fields interface{},
	actions []interface{},
	data map[string]interface{}) interface{} {

	return (&tabs.Component{}).Init().SetTabPanes(fields).SetTabBarExtraContent(extra)
}

// 详情页页面显示前回调
func (p *Template) BeforeDetailShowing(ctx *quark.Context, data map[string]interface{}) map[string]interface{} {
	return data
}
