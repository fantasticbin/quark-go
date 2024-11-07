package layout

import (
	"github.com/quarkcloudio/quark-go/v3"
)

// 布局模板
type Template struct {
	quark.Template
}

// 初始化
func (p *Template) Init(ctx *quark.Context) interface{} {
	return p
}

// 初始化模板
func (p *Template) TemplateInit(ctx *quark.Context) interface{} {
	return p
}

// 初始化路由映射
func (p *Template) RouteInit() interface{} {
	p.GET("/api/miniapp/layout/:resource/index", p.Render) // 渲染页面路由
	return p
}

// 头部
func (p *Template) Header(ctx *quark.Context) interface{} {
	return nil
}

// 页脚
func (p *Template) Footer(ctx *quark.Context) interface{} {
	return nil
}

// 表单数据
func (p *Template) Render(ctx *quark.Context) error {
	// 头部
	header := ctx.Template.(interface {
		Header(ctx *quark.Context) interface{}
	}).Header(ctx)

	// footer
	footer := ctx.Template.(interface {
		Footer(ctx *quark.Context) interface{}
	}).Footer(ctx)
	return ctx.JSON(200, map[string]interface{}{
		"header": header,
		"footer": footer,
	})
}
