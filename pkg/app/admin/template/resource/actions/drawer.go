package actions

import "github.com/quarkcloudio/quark-go/v3/pkg/builder"

type Drawer struct {
	Action
	Width          int  `json:"width"`          // 抽屉弹出层宽度
	DestroyOnClose bool `json:"destroyOnClose"` // 关闭时销毁弹出层里的子元素
}

// 初始化
func (p *Drawer) TemplateInit(ctx *builder.Context) interface{} {
	p.ActionType = "drawer"
	p.Width = 520

	return p
}

// 宽度
func (p *Drawer) GetWidth() int {
	return p.Width
}

// 关闭时销毁 Modal 里的子元素
func (p *Drawer) GetDestroyOnClose() bool {
	return p.DestroyOnClose
}

// 内容
func (p *Drawer) GetBody(ctx *builder.Context) interface{} {
	return nil
}

// 弹窗行为
func (p *Drawer) GetActions(ctx *builder.Context) []interface{} {
	return []interface{}{}
}
