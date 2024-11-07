package searches

import "github.com/quarkcloudio/quark-go/v3"

type Datetime struct {
	Search
}

// 初始化模板
func (p *Datetime) TemplateInit(ctx *quark.Context) interface{} {
	p.Component = "datetimeField"

	return p
}
