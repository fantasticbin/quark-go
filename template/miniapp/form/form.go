package form

import (
	"reflect"
	"strings"

	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/template/miniapp/component/form"
	"github.com/quarkcloudio/quark-go/v3/template/miniapp/page"
)

// 后台登录模板
type Template struct {
	page.Template
	FromStyle string
	Api       string
}

// 初始化
func (p *Template) Init(ctx *quark.Context) interface{} {
	return p
}

// 初始化模板
func (p *Template) TemplateInit(ctx *quark.Context) interface{} {

	// 初始化数据对象
	p.DB = db.Client

	// 标题
	p.Title = "QuarkGo"

	return p
}

// 初始化路由映射
func (p *Template) RouteInit() interface{} {
	p.GET("/api/miniapp/form/:resource/index", p.Render)  // 渲染页面路由
	p.Any("/api/miniapp/form/:resource/submit", p.Handle) // 表单提交路由

	return p
}

// 表单
func (p *Template) Form(api string, items []interface{}) *form.Component {
	return (&form.Component{}).
		Init().
		SetApi(api).
		SetBody(items)
}

// 表单项
func (p *Template) Field() *form.Field {
	return (&form.Field{})
}

// 表单项
func (p *Template) FormItem() *form.Field {
	return (&form.Field{})
}

// 字段
func (p *Template) Fields(ctx *quark.Context) []interface{} {
	return nil
}

// 行为
func (p *Template) Actions(ctx *quark.Context) []interface{} {
	return []interface{}{
		p.Action("提交", "primary").
			SetActionType("submit").
			SetFormType("button").
			SetBlock(true),
	}
}

// 表单数据
func (p *Template) Data(ctx *quark.Context) map[string]interface{} {
	return nil
}

// 内容
func (p *Template) Content(ctx *quark.Context) interface{} {

	fields := ctx.Template.(interface {
		Fields(ctx *quark.Context) []interface{}
	}).Fields(ctx)

	data := ctx.Template.(interface {
		Data(ctx *quark.Context) map[string]interface{}
	}).Data(ctx)

	actions := ctx.Template.(interface {
		Actions(ctx *quark.Context) []interface{}
	}).Actions(ctx)

	// 获取接口地址
	api := reflect.
		ValueOf(ctx.Template).
		Elem().
		FieldByName("Api").
		String()

	// 样式
	fromStyle := reflect.
		ValueOf(ctx.Template).
		Elem().
		FieldByName("FromStyle").
		String()

	if api == "" {
		api = "/api/miniapp/form/" + strings.ToLower(ctx.ResourceName()) + "/submit"
	}

	return p.Form(api, fields).
		SetStyle(fromStyle).
		SetModelValue(data).
		SetActions(actions)
}

// 执行表单
func (p *Template) Handle(ctx *quark.Context) error {
	return ctx.JSONError("请自行处理表单逻辑")
}
