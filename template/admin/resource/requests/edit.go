package requests

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/message"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
)

type EditRequest struct{}

// 表单数据
func (p *EditRequest) FillData(ctx *quark.Context) map[string]interface{} {
	result := map[string]interface{}{}
	id := ctx.Query("id", "")
	if id == "" {
		return result
	}

	// 模版实例
	template := ctx.Template.(types.Resourcer)

	// 模型结构体
	modelInstance := template.GetModel()

	// Gorm对象
	model := db.Client.Model(&modelInstance)

	// 创建编辑页查询
	query := template.BuildEditQuery(ctx, model)

	// 查询数据
	query.First(&result)

	// 获取字段
	updateFields := template.UpdateFields(ctx)

	// 给实例的Field属性赋值
	template.SetField(result)

	// 解析字段值
	fields := make(map[string]interface{})
	for _, field := range updateFields.([]interface{}) {

		// 字段名
		name := reflect.
			ValueOf(field).
			Elem().
			FieldByName("Name").
			String()

		if result[name] != nil {
			var fieldValue interface{}

			// 组件名称
			component := reflect.
				ValueOf(field).
				Elem().
				FieldByName("Component").
				String()

			if component == "datetimeField" {
				format := reflect.
					ValueOf(field).
					Elem().
					FieldByName("Format").
					String()

				format = strings.Replace(format, "YYYY", "2006", -1)
				format = strings.Replace(format, "MM", "01", -1)
				format = strings.Replace(format, "DD", "02", -1)
				format = strings.Replace(format, "HH", "15", -1)
				format = strings.Replace(format, "mm", "04", -1)
				format = strings.Replace(format, "ss", "05", -1)

				fieldValue = result[name].(time.Time).Format(format)
			} else if component == "dateField" {
				format := reflect.
					ValueOf(field).
					Elem().
					FieldByName("Format").
					String()

				format = strings.Replace(format, "YYYY", "2006", -1)
				format = strings.Replace(format, "MM", "01", -1)
				format = strings.Replace(format, "DD", "02", -1)

				fieldValue = result[name].(time.Time).Format(format)
			} else {
				fieldValue = result[name]
			}

			fields[name] = fieldValue
		}
	}

	return fields
}

// 获取表单初始化数据
func (p *EditRequest) Values(ctx *quark.Context) error {

	// 模版实例
	template := ctx.Template.(types.Resourcer)

	// 获取赋值数据
	data := p.FillData(ctx)

	// 获取初始数据
	data = template.BeforeEditing(ctx, data)

	// 解析数据
	for k, v := range data {
		getV, ok := v.(string)
		if ok {
			if strings.Contains(getV, "[") {
				var m []interface{}
				err := json.Unmarshal([]byte(getV), &m)
				if err == nil {
					v = m
				} else {
					if strings.Contains(getV, "{") {
						var m map[string]interface{}
						err := json.Unmarshal([]byte(getV), &m)
						if err == nil {
							v = m
						}
					}
				}
			} else {
				if strings.Contains(getV, "{") {
					var m map[string]interface{}
					err := json.Unmarshal([]byte(getV), &m)
					if err == nil {
						v = m
					}
				}
			}
		}

		data[k] = v
	}

	return ctx.JSON(200, message.Success("获取成功", "", data))
}
