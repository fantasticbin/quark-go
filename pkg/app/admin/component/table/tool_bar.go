package table

import "github.com/quarkcloudio/quark-go/v3/pkg/app/admin/component/component"

type ToolBar struct {
	component.Element
	Title        string      `json:"title"`
	SubTitle     string      `json:"subTitle"`
	Description  string      `json:"description"`
	Search       interface{} `json:"search"`
	Actions      interface{} `json:"actions"`
	Filter       interface{} `json:"filter"`
	MultipleLine bool        `json:"multipleLine"`
	Menu         interface{} `json:"menu"`
	Tabs         interface{} `json:"tabs"`
}

// 初始化
func (p *ToolBar) Init() *ToolBar {
	p.Component = "toolBar"

	p.SetKey("toolBar", component.DEFAULT_CRYPT)

	return p
}

// Set style.
func (p *ToolBar) SetStyle(style map[string]interface{}) *ToolBar {
	p.Style = style

	return p
}

// 标题
func (p *ToolBar) SetTitle(title string) *ToolBar {
	p.Title = title
	return p
}

// 子标题
func (p *ToolBar) SetSubTitle(subTitle string) *ToolBar {
	p.SubTitle = subTitle
	return p
}

// 描述
func (p *ToolBar) SetDescription(description string) *ToolBar {
	p.Description = description

	return p
}

// 查询区
func (p *ToolBar) SetSearch(search interface{}) *ToolBar {
	p.Search = search

	return p
}

// 操作区
func (p *ToolBar) SetAction(callback interface{}) *ToolBar {
	//  p.actions = callback(p.action);

	return p
}

// 批量设置操作区
func (p *ToolBar) SetActions(actions interface{}) *ToolBar {
	p.Actions = actions

	return p
}

// 过滤区，通常配合 LightFilter 使用
func (p *ToolBar) SetFilter(filter interface{}) *ToolBar {
	p.Filter = filter

	return p
}

// 是否多行展示
func (p *ToolBar) SetMultipleLine(multipleLine bool) *ToolBar {
	p.MultipleLine = multipleLine

	return p
}

// 菜单配置
func (p *ToolBar) SetMenu(menu interface{}) *ToolBar {
	p.Menu = menu

	return p
}

// 标签页配置，仅当 multipleLine 为 true 时有效
func (p *ToolBar) SetTabs(tabs interface{}) *ToolBar {
	p.Tabs = tabs

	return p
}

// 组件json序列化
func (p *ToolBar) JsonSerialize() *ToolBar {
	p.Component = "toolBar"

	return p
}
