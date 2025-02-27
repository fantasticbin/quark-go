package actions

import (
	"strings"

	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/actions"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
)

type CreateLinkAction struct {
	actions.Link
}

// 创建-跳转类型
func CreateLink() *CreateLinkAction {
	return &CreateLinkAction{}
}

// 初始化
func (p *CreateLinkAction) Init(ctx *quark.Context) interface{} {
	template := ctx.Template.(types.Resourcer)

	// 文字
	p.Name = "创建" + template.GetTitle()

	// 类型
	p.Type = "primary"

	// 图标
	p.Icon = "plus-circle"

	// 设置展示位置
	p.SetOnlyOnIndex(true)

	return p
}

// 跳转链接
func (p *CreateLinkAction) GetHref(ctx *quark.Context) string {
	return "#/layout/index?api=" + strings.Replace(ctx.Path(), "/index", "/create", -1)
}
