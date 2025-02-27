package resource

import (
	"reflect"
	"strings"

	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/descriptions"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form/fields/cascader"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form/fields/checkbox"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form/fields/radio"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form/fields/selectfield"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form/fields/treeselect"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/form/fields/when"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/table"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
)

// 列表字段
func (p *Template) IndexFields(ctx *quark.Context) interface{} {
	var items []interface{}

	fields := p.getFields(ctx)
	for _, v := range fields.([]interface{}) {
		if v, ok := v.(interface{ IsShownOnIndex() bool }); ok {
			if v.IsShownOnIndex() {
				items = append(items, v)
			}
		}
	}

	return items
}

// 表格列
func (p *Template) IndexTableColumns(ctx *quark.Context) interface{} {
	var columns []interface{}

	fields := p.IndexFields(ctx)
	for _, v := range fields.([]interface{}) {
		getColumn := p.fieldToColumn(ctx, v)
		if getColumn != nil {
			columns = append(columns, getColumn)
		}
	}

	// 资源实例
	template := ctx.Template.(types.Resourcer)

	// 行内行为
	indexTableRowActions := template.IndexTableRowActions(ctx)
	if len(indexTableRowActions.([]interface{})) > 0 {

		// 行为列标题
		columnTitle := p.GetTableActionColumnTitle()

		// 行为列宽度
		columnWidth := p.GetTableActionColumnWidth()

		// 构建列
		column := template.GetTableColumn(ctx).
			SetTitle(columnTitle).
			SetWidth(columnWidth).
			SetAttribute("action").
			SetValueType("option").
			SetActions(indexTableRowActions).
			SetFixed("right")
		columns = append(columns, column)
	}

	return columns
}

// 将表单项转换为表格列
func (p *Template) fieldToColumn(ctx *quark.Context, field interface{}) interface{} {
	var options interface{}

	reflectElem := reflect.
		ValueOf(field).
		Elem()

	// 是否可编辑
	columnField := reflectElem.
		FieldByName("Column").
		Interface()

	// 字段
	name := reflectElem.
		FieldByName("Name").
		String()

	// 文字
	label := reflectElem.
		FieldByName("Label").
		String()

	// 组件类型
	component := reflectElem.
		FieldByName("Component").
		String()

	// 列的对齐方式,left | right | center，只在列表页、详情页中有效
	align := reflectElem.
		FieldByName("Align").
		String()

	// （IE 下无效）列是否固定，可选 true (等效于 left) left rightr
	fixed := reflectElem.
		FieldByName("Fixed").
		Interface()

	// 是否可编辑
	editable := reflectElem.
		FieldByName("Editable").
		Bool()

	// 是否自动缩略
	ellipsis := reflectElem.
		FieldByName("Ellipsis").
		Bool()

	// 是否支持复制
	copyable := reflectElem.
		FieldByName("Copyable").
		Bool()

	// 表头的筛选菜单项，当值为 true 时，自动使用 valueEnum 生成
	filters := reflectElem.
		FieldByName("Filters").
		Interface()

	// 查询表单中的权重，权重大排序靠前
	order := reflectElem.
		FieldByName("Order").
		Int()

	// 可排序列
	sorter := reflectElem.
		FieldByName("Sorter").
		Interface()

	// 包含列的数量
	span := reflectElem.
		FieldByName("Span").
		Int()

	// 设置列宽
	columnWidth := reflectElem.
		FieldByName("ColumnWidth").
		Int()

	column := columnField.(*table.Column).
		SetTitle(label).
		SetAttribute(name).
		SetAlign(align).
		SetFixed(fixed).
		SetEllipsis(ellipsis).
		SetCopyable(copyable).
		SetFilters(filters).
		SetOrder(int(order)).
		SetSorter(sorter).
		SetSpan(int(span)).
		SetWidth(int(columnWidth))

	switch component {
	case "idField":
		// 是否显示在列表
		onIndexDisplayed := reflectElem.
			FieldByName("OnIndexDisplayed").
			Bool()
		if onIndexDisplayed {
			column = column.SetValueType("text")
		} else {
			return nil
		}
	case "hiddenField":
		return nil
	case "textField":
		column = column.SetValueType("text")
	case "textAreaField":
		column = column.SetValueType("text")
	case "treeSelectField":

		// 获取属性
		options = field.(interface{ GetOptions() []treeselect.TreeData }).GetOptions()

		// 设置表格列
		column = column.
			SetValueType("treeSelect").
			SetFieldProps(map[string]interface{}{
				"options": options,
			})
	case "cascaderField":
		// 获取属性
		options = field.(interface{ GetOptions() []cascader.Option }).GetOptions()

		// 设置表格列
		column = column.
			SetValueType("cascader").
			SetFieldProps(map[string]interface{}{
				"options": options,
			})
	case "selectField":
		// 获取属性
		options = field.(interface{ GetOptions() []selectfield.Option }).GetOptions()

		// 设置表格列
		column = column.
			SetValueType("select").
			SetFieldProps(map[string]interface{}{
				"options": options,
			})

		// 是否设置了过滤项
		if getfilters, ok := filters.(bool); ok {
			if getfilters {
				// 获取值的枚举，会自动转化把值当成 key 来取出要显示的内容
				valueEnum := field.(interface{ GetValueEnum() interface{} }).GetValueEnum()

				// 设置值的枚举
				column.SetValueEnum(valueEnum)
			}
		}
	case "checkboxField":
		// 获取属性
		options = field.(interface{ GetOptions() []checkbox.Option }).GetOptions()

		// 设置表格列
		column = column.
			SetValueType("checkbox").
			SetFieldProps(map[string]interface{}{
				"options": options,
			})

		// 是否设置了过滤项
		if getfilters, ok := filters.(bool); ok {
			if getfilters {
				// 获取值的枚举，会自动转化把值当成 key 来取出要显示的内容
				valueEnum := field.(interface{ GetValueEnum() interface{} }).GetValueEnum()

				// 设置值的枚举
				column.SetValueEnum(valueEnum)
			}
		}
	case "radioField":
		// 获取属性
		options = field.(interface{ GetOptions() []radio.Option }).GetOptions()

		// 设置表格列
		column = column.
			SetValueType("radio").
			SetFieldProps(map[string]interface{}{
				"options": options,
			})

		// 是否设置了过滤项
		if getfilters, ok := filters.(bool); ok {
			if getfilters {
				// 获取值的枚举，会自动转化把值当成 key 来取出要显示的内容
				valueEnum := field.(interface{ GetValueEnum() interface{} }).GetValueEnum()

				// 设置值的枚举
				column.SetValueEnum(valueEnum)
			}
		}
	case "switchField":
		// 获取属性
		options = field.(interface{ GetOptions() interface{} }).GetOptions()

		// 设置表格列
		column = column.
			SetValueType("select").
			SetValueEnum(options)

		// 是否设置了过滤项
		if getfilters, ok := filters.(bool); ok {
			if getfilters {
				// 获取值的枚举，会自动转化把值当成 key 来取出要显示的内容
				valueEnum := field.(interface{ GetValueEnum() interface{} }).GetValueEnum()

				// 设置值的枚举
				column.SetValueEnum(valueEnum)
			}
		}
	case "imageField":
		column = column.SetValueType("image")
	default:
		column = column.SetValueType(component)
	}

	if editable {

		// 可编辑api地址
		editableApi := strings.Replace(ctx.Path(), "/index", "/editable", -1)

		// 设置编辑项
		column = column.SetEditable(component, options, editableApi)
	}

	return column
}

// 创建页字段
func (p *Template) CreationFields(ctx *quark.Context) interface{} {
	var items []interface{}

	fields := p.getFields(ctx)
	for _, v := range fields.([]interface{}) {
		if v, ok := v.(interface{ IsShownOnCreation() bool }); ok {
			if v.IsShownOnCreation() {
				items = append(items, v)
			}
		}
	}

	return items
}

// 不包含When组件内字段的创建页字段
func (p *Template) CreationFieldsWithoutWhen(ctx *quark.Context) interface{} {
	var items []interface{}

	fields := p.getFieldsWithoutWhen(ctx)
	for _, v := range fields.([]interface{}) {
		if v, ok := v.(interface{ IsShownOnCreation() bool }); ok {
			if v.IsShownOnCreation() {
				items = append(items, v)
			}
		}
	}

	return items
}

// 包裹在组件内的创建页字段
func (p *Template) CreationFieldsWithinComponents(ctx *quark.Context) interface{} {

	// 资源实例
	template := ctx.Template.(types.Resourcer)

	// 获取字段
	fields := template.Fields(ctx)

	// 解析创建页表单组件内的字段
	items := p.CreationFormFieldsParser(ctx, fields)

	return items
}

// 解析创建页表单组件内的字段
func (p *Template) CreationFormFieldsParser(ctx *quark.Context, fields interface{}) interface{} {
	items := []interface{}{}

	// 解析字段
	if fields, ok := fields.([]interface{}); ok {
		for _, v := range fields {
			hasBody := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Body").
				IsValid()
			if hasBody {

				// 获取内容值
				body := reflect.
					ValueOf(v).
					Elem().
					FieldByName("Body").
					Interface()

				// 解析值
				getFields := p.CreationFormFieldsParser(ctx, body)

				// 更新值
				reflect.
					ValueOf(v).
					Elem().
					FieldByName("Body").
					Set(reflect.ValueOf(getFields))

				items = append(items, v)
			} else {
				component := reflect.
					ValueOf(v).
					Elem().
					FieldByName("Component").
					String()
				if strings.Contains(component, "Field") {
					whenIsValid := reflect.
						ValueOf(v).
						Elem().
						FieldByName("When").
						IsValid()
					if whenIsValid {
						getWhen := v.(interface {
							GetWhen() *when.Component
						}).GetWhen()
						if getWhen != nil {
							whenItems := getWhen.Items
							for _, sv := range whenItems {
								if sv.Body != nil {
									if svBody, ok := sv.Body.([]interface{}); ok {
										// 解析值
										getSvBody := p.CreationFormFieldsParser(ctx, svBody)
										// 更新值
										reflect.
											ValueOf(sv).
											Elem().
											FieldByName("Body").
											Set(reflect.ValueOf(getSvBody))
									} else {
										// 判断是否在编辑页面
										if getSvBody, ok := sv.Body.(interface{ IsShownOnCreation() bool }); ok {
											if getSvBody.IsShownOnCreation() {

												// 生成前端验证规则
												getSvBody.(interface{ BuildFrontendRules(string) interface{} }).BuildFrontendRules(ctx.Path())

												// 更新值
												reflect.
													ValueOf(sv).
													Elem().
													FieldByName("Body").
													Set(reflect.ValueOf(getSvBody))
											}
										}
									}
								}
							}
						}
					}
					// 判断是否在创建页面
					if v, ok := v.(interface{ IsShownOnCreation() bool }); ok {
						if v.IsShownOnCreation() {

							// 生成前端验证规则
							v.(interface{ BuildFrontendRules(string) interface{} }).BuildFrontendRules(ctx.Path())

							// 组合数据
							items = append(items, v)
						}
					}
				} else {
					items = append(items, v)
				}
			}
		}
	}

	return items
}

// 编辑页字段
func (p *Template) UpdateFields(ctx *quark.Context) interface{} {
	var items []interface{}

	fields := p.getFields(ctx)
	for _, v := range fields.([]interface{}) {
		if v, ok := v.(interface{ IsShownOnUpdate() bool }); ok {
			if v.IsShownOnUpdate() {
				items = append(items, v)
			}
		}
	}

	return items
}

// 不包含When组件内字段的编辑页字段
func (p *Template) UpdateFieldsWithoutWhen(ctx *quark.Context) interface{} {
	var items []interface{}

	fields := p.getFieldsWithoutWhen(ctx)
	for _, v := range fields.([]interface{}) {
		if v, ok := v.(interface{ IsShownOnUpdate() bool }); ok {
			if v.IsShownOnUpdate() {
				items = append(items, v)
			}
		}
	}

	return items
}

// 包裹在组件内的编辑页字段
func (p *Template) UpdateFieldsWithinComponents(ctx *quark.Context) interface{} {

	// 资源实例
	template := ctx.Template.(types.Resourcer)

	// 获取字段
	fields := template.Fields(ctx)

	// 解析编辑页表单组件内的字段
	items := p.UpdateFormFieldsParser(ctx, fields)

	return items
}

// 解析编辑页表单组件内的字段
func (p *Template) UpdateFormFieldsParser(ctx *quark.Context, fields interface{}) interface{} {
	items := []interface{}{}

	// 解析字段
	if fields, ok := fields.([]interface{}); ok {
		for _, v := range fields {
			hasBody := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Body").
				IsValid()
			if hasBody {

				// 获取内容值
				body := reflect.
					ValueOf(v).
					Elem().
					FieldByName("Body").
					Interface()

				// 解析值
				getFields := p.UpdateFormFieldsParser(ctx, body)

				// 更新值
				reflect.
					ValueOf(v).
					Elem().
					FieldByName("Body").
					Set(reflect.ValueOf(getFields))

				items = append(items, v)
			} else {
				component := reflect.
					ValueOf(v).
					Elem().
					FieldByName("Component").
					String()
				if strings.Contains(component, "Field") {
					whenIsValid := reflect.
						ValueOf(v).
						Elem().
						FieldByName("When").
						IsValid()
					if whenIsValid {
						getWhen := v.(interface {
							GetWhen() *when.Component
						}).GetWhen()
						if getWhen != nil {
							whenItems := getWhen.Items
							for _, sv := range whenItems {
								if sv.Body != nil {
									if svBody, ok := sv.Body.([]interface{}); ok {
										// 解析值
										getSvBody := p.UpdateFormFieldsParser(ctx, svBody)
										// 更新值
										reflect.
											ValueOf(sv).
											Elem().
											FieldByName("Body").
											Set(reflect.ValueOf(getSvBody))
									} else {
										// 判断是否在编辑页面
										if getSvBody, ok := sv.Body.(interface{ IsShownOnUpdate() bool }); ok {
											if getSvBody.IsShownOnUpdate() {

												// 生成前端验证规则
												getSvBody.(interface{ BuildFrontendRules(string) interface{} }).BuildFrontendRules(ctx.Path())

												// 更新值
												reflect.
													ValueOf(sv).
													Elem().
													FieldByName("Body").
													Set(reflect.ValueOf(getSvBody))
											}
										}
									}
								}
							}
						}
					}
					// 判断是否在编辑页面
					if v, ok := v.(interface{ IsShownOnUpdate() bool }); ok {
						if v.IsShownOnUpdate() {

							// 生成前端验证规则
							v.(interface{ BuildFrontendRules(string) interface{} }).BuildFrontendRules(ctx.Path())

							// 组合数据
							items = append(items, v)
						}
					}
				} else {
					items = append(items, v)
				}
			}
		}
	}

	return items
}

// 解析表单组件内的字段
func (p *Template) FormFieldsParser(ctx *quark.Context, fields interface{}) interface{} {
	if fields, ok := fields.([]interface{}); ok {
		for k, v := range fields {
			hasBody := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Body").
				IsValid()
			if hasBody {

				// 获取内容值
				body := reflect.
					ValueOf(v).
					Elem().
					FieldByName("Body").
					Interface()

				// 解析值
				getFields := p.FormFieldsParser(ctx, body)

				// 更新值
				reflect.
					ValueOf(v).
					Elem().
					FieldByName("Body").
					Set(reflect.ValueOf(getFields))
			} else {
				component := reflect.
					ValueOf(v).
					Elem().
					FieldByName("Component").
					String()
				if strings.Contains(component, "Field") {
					whenIsValid := reflect.
						ValueOf(v).
						Elem().
						FieldByName("When").
						IsValid()
					if whenIsValid {
						getWhen := v.(interface {
							GetWhen() *when.Component
						}).GetWhen()
						if getWhen != nil {
							whenItems := getWhen.Items
							for _, sv := range whenItems {
								if sv.Body != nil {
									if svBody, ok := sv.Body.([]interface{}); ok {
										// 解析值
										getSvBody := p.FormFieldsParser(ctx, svBody)
										// 更新值
										reflect.
											ValueOf(sv).
											Elem().
											FieldByName("Body").
											Set(reflect.ValueOf(getSvBody))
									} else {
										// 生成前端验证规则
										sv.Body.(interface{ BuildFrontendRules(string) interface{} }).BuildFrontendRules(ctx.Path())
										// 更新值
										reflect.
											ValueOf(sv).
											Elem().
											FieldByName("Body").
											Set(reflect.ValueOf(sv.Body))
									}
								}
							}
						}
					}

					// 生成前端验证规则
					v.(interface{ BuildFrontendRules(string) interface{} }).BuildFrontendRules(ctx.Path())
				}
			}

			fields[k] = v
		}
	}

	return fields
}

// 详情页字段
func (p *Template) DetailFields(ctx *quark.Context) interface{} {
	var items []interface{}

	fields := p.getFields(ctx)
	for _, v := range fields.([]interface{}) {
		if v, ok := v.(interface{ IsShownOnDetail() bool }); ok {
			if v.IsShownOnDetail() {
				items = append(items, v)
			}
		}
	}

	return items
}

// 包裹在组件内的详情页字段
func (p *Template) DetailFieldsWithinComponents(ctx *quark.Context, data map[string]interface{}) interface{} {
	var (
		items         []interface{}
		componentType = "description"
	)

	// 资源实例
	template := ctx.Template.(types.Resourcer)

	// 解析字段
	fields := template.Fields(ctx)
	for _, v := range fields {

		hasBody := reflect.
			ValueOf(v).
			Elem().
			FieldByName("Body").
			IsValid()

		// 解析body数据
		if hasBody {
			body := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Body").
				Interface()

			var subItems []interface{}
			for _, sv := range body.([]interface{}) {
				if sv, ok := sv.(interface{ IsShownOnDetail() bool }); ok {
					if sv.IsShownOnDetail() {
						getColumn := p.fieldToColumn(ctx, sv)
						subItems = append(subItems, getColumn)
					}
				}
			}

			descriptions := (&descriptions.Component{}).
				Init().
				SetStyle(map[string]interface{}{
					"padding": "24px",
				}).
				SetTitle("").
				SetColumn(2).
				SetColumns(subItems).
				SetDataSource(data).
				SetActions(p.DetailActions(ctx))

			v.(interface{ SetBody(interface{}) interface{} }).SetBody(descriptions)
			items = append(items, v)
		} else {
			if v, ok := v.(interface{ IsShownOnDetail() bool }); ok {
				if v.IsShownOnDetail() {
					getColumn := p.fieldToColumn(ctx, v)
					if getColumn != nil {
						items = append(items, getColumn)
					}
				}
			}
		}
	}

	if componentType == "description" {
		return (&descriptions.Component{}).
			Init().
			SetStyle(map[string]interface{}{
				"padding": "24px",
			}).
			SetTitle("").
			SetColumn(2).
			SetColumns(items).
			SetDataSource(data).
			SetActions(p.DetailActions(ctx))
	}

	return items
}

// 导出字段
func (p *Template) ExportFields(ctx *quark.Context) interface{} {
	var items []interface{}

	fields := p.getFields(ctx)
	for _, v := range fields.([]interface{}) {
		if v, ok := v.(interface{ IsShownOnExport() bool }); ok {
			if v.IsShownOnExport() {
				items = append(items, v)
			}
		}
	}

	return items
}

// 导入字段
func (p *Template) ImportFields(ctx *quark.Context) interface{} {
	var items []interface{}

	fields := p.getFields(ctx)
	for _, v := range fields.([]interface{}) {
		if v, ok := v.(interface{ IsShownOnImport() bool }); ok {
			if v.IsShownOnImport() {
				items = append(items, v)
			}
		}
	}

	return items
}

// 不包含When组件内字段的导入字段
func (p *Template) ImportFieldsWithoutWhen(ctx *quark.Context) interface{} {
	var items []interface{}

	fields := p.getFieldsWithoutWhen(ctx)
	for _, v := range fields.([]interface{}) {
		if v, ok := v.(interface{ IsShownOnImport() bool }); ok {
			if v.IsShownOnImport() {
				items = append(items, v)
			}
		}
	}

	return items
}

// 获取字段
func (p *Template) getFields(ctx *quark.Context) interface{} {

	// 资源实例
	template := ctx.Template.(types.Resourcer)

	// 获取字段
	fields := template.Fields(ctx)

	return p.findFields(fields, true)
}

// 获取不包含When组件的字段
func (p *Template) getFieldsWithoutWhen(ctx *quark.Context) interface{} {

	// 资源实例
	template := ctx.Template.(types.Resourcer)

	// 获取字段
	fields := template.Fields(ctx)

	return p.findFields(fields, false)
}

// 查找字段
func (p *Template) findFields(fields interface{}, when bool) interface{} {
	var items []interface{}

	if fields, ok := fields.([]interface{}); ok {
		for _, v := range fields {
			hasBody := reflect.
				ValueOf(v).
				Elem().
				FieldByName("Body").
				IsValid()
			if hasBody {
				body := reflect.
					ValueOf(v).
					Elem().
					FieldByName("Body").
					Interface()
				getItems := p.findFields(body, true)
				if getItems, ok := getItems.([]interface{}); ok {
					if len(getItems) > 0 {
						items = append(items, getItems...)
					}
				}
			} else {
				component := reflect.
					ValueOf(v).
					Elem().
					FieldByName("Component").
					String()
				if strings.Contains(component, "Field") {
					items = append(items, v)
					if when {
						whenFields := p.getWhenFields(v)
						if len(whenFields) > 0 {
							items = append(items, whenFields...)
						}
					}
				}
			}
		}
	}

	return items
}

// 获取When组件中的字段
func (p *Template) getWhenFields(item interface{}) []interface{} {
	var items []interface{}
	whenIsValid := reflect.
		ValueOf(item).
		Elem().
		FieldByName("When").
		IsValid()
	if !whenIsValid {
		return items
	}

	getWhen := item.(interface {
		GetWhen() *when.Component
	}).GetWhen()

	if getWhen == nil {
		return items
	}
	whenItems := getWhen.Items
	if whenItems == nil {
		return items
	}

	for _, v := range whenItems {
		if v.Body != nil {
			if body, ok := v.Body.([]interface{}); ok {
				if len(body) > 0 {
					items = append(items, body...)
				}
			} else {
				items = append(items, v.Body)
			}
		}
	}

	return items
}
