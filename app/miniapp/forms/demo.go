package forms

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/miniapp/form"
)

type Demo struct {
	form.Template
}

// 初始化
func (p *Demo) Init(ctx *quark.Context) interface{} {

	return p
}

// 字段
func (p *Demo) Fields(ctx *quark.Context) []interface{} {
	return []interface{}{
		p.Field().Input("username", "姓名"),
		p.Field().Input("password", "密码"),
	}
}
