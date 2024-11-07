package pages

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/miniapp/component/col"
	"github.com/quarkcloudio/quark-go/v3/template/miniapp/page"
)

type My struct {
	page.Template
}

// 初始化
func (p *My) Init(ctx *quark.Context) interface{} {
	return p
}

// 组件渲染
func (p *My) Content(ctx *quark.Context) interface{} {
	return []interface{}{
		p.Row([]*col.Component{
			p.Col(24, "我的"),
		}).SetStyle("text-align:center;"),
	}
}
