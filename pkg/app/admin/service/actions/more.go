package actions

import (
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/template/resource/actions"
	"github.com/quarkcloudio/quark-go/v3/pkg/builder"
)

type MoreAction struct {
	actions.Dropdown
}

// 更多，More() | More("更多") | More("更多", []interface{}{})
func More(options ...interface{}) *MoreAction {
	moreAction := &MoreAction{}

	moreAction.Name = "更多"
	if len(options) == 1 {
		moreAction.Name = options[0].(string)
	}

	if len(options) == 2 {
		moreAction.Name = options[0].(string)
		moreAction.Actions = options[1].([]interface{})
	}

	return moreAction
}

// 初始化
func (p *MoreAction) Init(ctx *builder.Context) interface{} {

	// 下拉框箭头是否显示
	p.Arrow = true

	// 菜单弹出位置：bottomLeft bottomCenter bottomRight topLeft topCenter topRight
	p.Placement = "bottomLeft"

	// 触发下拉的行为, 移动端不支持 hover,Array<click|hover|contextMenu>
	p.Trigger = []string{"hover"}

	// 下拉根元素的样式
	p.OverlayStyle = map[string]interface{}{
		"zIndex": 999,
	}

	// 设置按钮类型,primary | ghost | dashed | link | text | default
	p.Type = "link"

	// 设置按钮大小,large | middle | small | default
	p.Size = "small"

	// 设置展示位置
	p.SetOnlyOnIndexTableRow(true)

	return p
}

// 下拉菜单行为
func (p *MoreAction) SetActions(actions []interface{}) interface{} {
	p.Actions = actions

	return p
}
