package pages

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/miniapp/component/col"
	"github.com/quarkcloudio/quark-go/v3/template/miniapp/page"
)

type Index struct {
	page.Template
}

// 初始化
func (p *Index) Init(ctx *quark.Context) interface{} {
	return p
}

// 组件渲染
func (p *Index) Content(ctx *quark.Context) interface{} {
	return []interface{}{
		p.Row([]*col.Component{
			p.Col(24, "Hello World!"),
		}).SetStyle("text-align:center;"),
	}
}
