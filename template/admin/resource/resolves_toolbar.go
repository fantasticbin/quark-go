package resource

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
)

// 列表工具栏
func (p *Template) IndexTableMenus(ctx *quark.Context) interface{} {

	// 模版实例
	template := ctx.Template.(types.Resourcer)

	menus := template.Menus(ctx)

	return menus
}
