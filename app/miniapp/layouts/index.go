package layouts

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/miniapp/component/tabbar"
	"github.com/quarkcloudio/quark-go/v3/template/miniapp/layout"
)

type Index struct {
	layout.Template
}

// 页脚
func (p *Index) Footer(ctx *quark.Context) interface{} {
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
