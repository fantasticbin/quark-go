package login

import "github.com/quarkcloudio/quark-go/v3/pkg/builder"

type Loginer interface {

	// 模版接口
	builder.Templater

	// 获取登录接口
	GetApi() string

	// 获取登录成功后跳转地址
	GetRedirect() string

	// 获取登录页面Logo
	GetLogo() interface{}

	// 获取登录页面标题
	GetTitle() string

	// 获取登录页面子标题
	GetSubTitle() string

	// 验证码ID
	CaptchaId(ctx *builder.Context) error

	// 生成验证码
	Captcha(ctx *builder.Context) error

	// 字段
	Fields(ctx *builder.Context) []interface{}

	// 登录方法
	Handle(ctx *builder.Context) error

	// 退出方法
	Logout(ctx *builder.Context) error

	// 包裹在组件内的创建页字段
	FieldsWithinComponents(ctx *builder.Context) interface{}

	// 组件渲染
	Render(ctx *builder.Context) error
}
