package layouts

import (
	"github.com/quarkcloudio/quark-go/v3/pkg/app/miniapp/component/tabbar"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/miniapp/template/layout"
	"github.com/quarkcloudio/quark-go/v3/pkg/builder"
)

type Index struct {
	layout.Template
}

// 页脚
func (p *Index) Footer(ctx *builder.Context) interface{} {
	return tabbar.New().
		SetBottom(true).
		SetItems([]*tabbar.Item{
			tabbar.NewItem().
				SetIcon("home").
				SetTabTitle("首页").
				SetTo("/pages/index"),
			tabbar.NewItem().
				SetIcon("my").
				SetTabTitle("我的").
				SetTo("/pages/my"),
		})
}
