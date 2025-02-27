package actions

import "github.com/quarkcloudio/quark-go/v3"

type DrawerForm struct {
	Action
	Width          int    `json:"width"`          // 抽屉弹出层宽度
	DestroyOnClose bool   `json:"destroyOnClose"` // 关闭时销毁弹出层里的子元素
	CancelText     string `json:"cancelText"`     // 获取取消按钮文案
	SubmitText     string `json:"submitText"`     // 获取提交按钮文案
}

// 初始化
func (p *DrawerForm) TemplateInit(ctx *quark.Context) interface{} {
	p.ActionType = "drawerForm"
	p.Width = 520
	p.Reload = "table"
	p.CancelText = "取消"
	p.SubmitText = "提交"

	return p
}

// 表单字段
func (p *DrawerForm) Fields(ctx *quark.Context) []interface{} {
	return []interface{}{}
}

// 表单数据（异步获取）
func (p *DrawerForm) Data(ctx *quark.Context) map[string]interface{} {
	return map[string]interface{}{}
}

// 宽度
func (p *DrawerForm) GetWidth() int {
	return p.Width
}

// 关闭时销毁 Modal 里的子元素
func (p *DrawerForm) GetDestroyOnClose() bool {
	return p.DestroyOnClose
}

// 获取取消按钮文案
func (p *DrawerForm) GetCancelText() string {
	return p.CancelText
}

// 获取提交按钮文案
func (p *DrawerForm) GetSubmitText() string {
	return p.SubmitText
}
