package layouts

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/layout"
)

type Index struct {
	layout.Template
}

// 初始化
func (p *Index) Init(ctx *quark.Context) interface{} {
	return p
}
