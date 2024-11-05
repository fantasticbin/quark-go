package pages

import (
	"github.com/quarkcloudio/quark-go/v3/pkg/app/miniapp/component/col"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/miniapp/template/page"
	"github.com/quarkcloudio/quark-go/v3/pkg/builder"
)

type My struct {
	page.Template
}

// 初始化
func (p *My) Init(ctx *builder.Context) interface{} {
	return p
}

// 组件渲染
func (p *My) Content(ctx *builder.Context) interface{} {
	return []interface{}{
		p.Row([]*col.Component{
			p.Col(24, "我的"),
		}).SetStyle("text-align:center;"),
	}
}
