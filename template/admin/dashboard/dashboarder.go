package dashboard

import (
	"github.com/quarkcloudio/quark-go/v3"
)

type Dashboarder interface {

	// 模版接口
	quark.Templater

	// 获取页面标题
	GetTitle() string

	// 获取页面子标题
	GetSubTitle() string

	// 页面是否携带返回Icon
	GetBackIcon() bool

	// 内容
	Cards(ctx *quark.Context) []interface{}

	// 页面组件渲染
	PageComponentRender(ctx *quark.Context, body interface{}) interface{}

	// 页面容器组件渲染
	PageContainerComponentRender(ctx *quark.Context, body interface{}) interface{}

	// 组件渲染
	Render(ctx *quark.Context) error
}
