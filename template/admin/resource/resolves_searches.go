package resource

import (
	"strings"

	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form/fields/radio"

	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form/fields/cascader"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form/fields/selectfield"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form/fields/treeselect"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
)

// 列表页搜索表单
func (p *Template) IndexSearches(ctx *quark.Context) interface{} {

	// 模版实例
	template := ctx.Template.(types.Resourcer)

	// 搜索项
	searches := template.Searches(ctx)

	// 搜索组件
	search := template.GetTableSearch(ctx)

	// 是否携带导出功能
	withExport := template.GetWithExport()
	if withExport {
		search = search.
			SetExportText("导出").
			SetExportApi(strings.Replace(ExportPath, ":resource", ctx.Param("resource"), -1))
	}

	// 解析搜索项
	for _, v := range searches {

		// 搜索栏表单项
		var item interface{}
		var field = &Field{}

		// 搜索实例
		searchInstance := v.(types.Searcher)

		// 初始化模版
		searchInstance.TemplateInit(ctx)

		// 初始化
		searchInstance.Init(ctx)

		// 获取组件名称
		component := searchInstance.GetComponent()

		// label 标签的文本
		label := searchInstance.GetName()

		// 字段名，支持数组
		name := searchInstance.GetColumn(v)

		// 获取接口
		api := searchInstance.GetApi()

		// 获取属性
		options := searchInstance.Options(ctx)

		// 获取Select组件的Load
		load := searchInstance.Load(ctx)

		// 构建组件
		switch component {
		case "textField":
			item = field.
				Text(name, label).
				SetWidth(nil)
		case "selectField":
			item = field.
				Select(name, label).
				SetWidth(nil).
				SetOptions(options.([]selectfield.Option)).
				SetLoad(load["field"], load["api"])
		case "radioField":
			item = field.
				Radio(name, label).
				SetOptions(options.([]radio.Option)).
				SetOptionType("button").
				SetButtonStyle("solid")
		case "multipleSelectField":
			item = field.
				Select(name, label).
				SetMode("multiple").
				SetWidth(nil).
				SetOptions(options.([]selectfield.Option))
		case "dateField":
			item = field.
				Date(name, label).
				SetWidth(nil)
		case "datetimeField":
			item = field.
				Datetime(name, label).
				SetWidth(nil)
		case "dateRangeField":
			item = field.
				DateRange(name, label).
				SetWidth(nil)
		case "datetimeRangeField":
			item = field.
				DatetimeRange(name, label).
				SetWidth(nil)
		case "cascaderField":
			item = field.
				Cascader(name, label).
				SetApi(api).
				SetWidth(nil).
				SetOptions(options.([]cascader.Option))
		case "treeSelectField":
			item = field.
				TreeSelect(name, label).
				SetWidth(nil).
				SetTreeData(options.([]treeselect.TreeData))
		}

		search = search.SetItems(item)
	}

	return search
}
