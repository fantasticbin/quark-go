package actions

import "github.com/quarkcloudio/quark-go/v3"

type Link struct {
	Action
	Href   string `json:"href"`   // 获取跳转链接
	Target string `json:"target"` // 相当于 a 链接的 target 属性，href 存在时生效，_blank | _self | _parent | _top
}

// 初始化
func (p *Link) TemplateInit(ctx *quark.Context) interface{} {
	p.ActionType = "link"
	p.Target = "_self"

	return p
}

// 获取跳转链接
func (p *Link) GetHref(ctx *quark.Context) string {
	return p.Href
}

// 相当于 a 链接的 target 属性，href 存在时生效
func (p *Link) GetTarget(ctx *quark.Context) string {
	return p.Target
}
