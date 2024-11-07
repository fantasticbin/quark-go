package table

import "github.com/quarkcloudio/quark-go/v3/template/admin/component/component"

type Search struct {
	component.Element
	FilterType       string        `json:"filterType,omitempty"`
	SearchText       string        `json:"searchText,omitempty"`
	ResetText        string        `json:"resetText,omitempty"`
	SubmitText       string        `json:"submitText,omitempty"`
	LabelWidth       int           `json:"labelWidth,omitempty"`
	Span             int           `json:"span,omitempty"`
	ClassName        string        `json:"className,omitempty"`
	DefaultCollapsed bool          `json:"defaultCollapsed,omitempty"`
	ShowHiddenNum    bool          `json:"showHiddenNum,omitempty"`
	ExportText       string        `json:"exportText,omitempty"`
	ExportApi        string        `json:"exportApi,omitempty"`
	Items            []interface{} `json:"items"`
}

// 初始化
func (p *Search) Init() *Search {
	p.Component = "search"
	p.DefaultCollapsed = true
	p.ResetText = "重置"
	p.SearchText = "查询"

	p.SetKey(component.DEFAULT_KEY, component.DEFAULT_CRYPT)

	return p
}

// Set style.
func (p *Search) SetStyle(style map[string]interface{}) *Search {
	p.Style = style

	return p
}

// 过滤表单类型，'query' | 'light'
func (p *Search) SetFilterType(filterType string) *Search {
	p.FilterType = filterType

	return p
}

// 默认状态下是否折叠超出的表单项
func (p *Search) SetCollapsed(collapsed bool) *Search {
	p.DefaultCollapsed = collapsed

	return p
}

// 查询按钮的文本
func (p *Search) SetSearchText(searchText string) *Search {
	p.SearchText = searchText

	return p
}

// 重置按钮的文本
func (p *Search) SetResetText(resetText string) *Search {
	p.ResetText = resetText

	return p
}

// 提交按钮的文本
func (p *Search) SetSubmitText(submitText string) *Search {
	p.SubmitText = submitText

	return p
}

// 封装的搜索 Form 的 className
func (p *Search) SetClassName(className string) *Search {
	p.ClassName = className

	return p
}

// label 宽度,number | 'auto'
func (p *Search) SetLabelWidth(labelWidth int) *Search {
	p.LabelWidth = labelWidth

	return p
}

// 表单项宽度,number[0 - 24]
func (p *Search) SetSpan(span int) *Search {
	p.Span = span

	return p
}

// 是否显示收起之后显示隐藏个数
func (p *Search) SetSplit(showHiddenNum bool) *Search {
	p.ShowHiddenNum = showHiddenNum

	return p
}

// 导出按钮文字
func (p *Search) SetExportText(exportText string) *Search {
	p.ExportText = exportText

	return p
}

// 导出数据接口
func (p *Search) SetExportApi(exportApi string) *Search {
	p.ExportApi = exportApi

	return p
}

// 设置搜索表单项
func (p *Search) SetItems(item interface{}) *Search {
	p.Items = append(p.Items, item)

	return p
}

// 组件json序列化
func (p *Search) JsonSerialize() *Search {
	p.Component = "search"

	return p
}
