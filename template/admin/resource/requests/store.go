package requests

import (
	"encoding/json"
	"reflect"

	"github.com/gobeam/stringy"
	"github.com/gookit/goutil/structs"
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/message"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
)

type StoreRequest struct{}

// 执行行为
func (p *StoreRequest) Handle(ctx *quark.Context, data map[string]interface{}) error {

	// 模版实例
	template := ctx.Template.(types.Resourcer)

	// 模型结构体
	modelInstance := template.GetModel()

	// 数据实例
	dataInstance := template.GetModel()

	// 验证数据合法性
	validator := template.ValidatorForCreation(ctx, data)
	if validator != nil {
		return ctx.JSON(200, message.Error(validator.Error()))
	}

	// 保存前回调
	data, err := template.BeforeSaving(ctx, data)
	if err != nil {
		return ctx.JSON(200, message.Error(err.Error()))
	}

	// 重组数据
	newData := map[string]interface{}{}
	for k, v := range data {
		nv := v

		// 将数组、map数据转换为字符串存储
		if gv, ok := v.([]interface{}); ok {
			nv, _ = json.Marshal(gv)
		}
		if gv, ok := v.([]map[string]interface{}); ok {
			nv, _ = json.Marshal(gv)
		}
		if gv, ok := v.(map[string]interface{}); ok {
			nv, _ = json.Marshal(gv)
		}

		camelCaseName := stringy.
			New(k).
			CamelCase("?", "")

		fieldIsValid := reflect.
			ValueOf(modelInstance).
			Elem().
			FieldByName(camelCaseName).
			IsValid()
		if fieldIsValid {
			newData[k] = nv
		}
	}

	// 结构体赋值
	structs.SetValues(dataInstance, newData)

	// 获取对象
	model := db.Client.Model(modelInstance).Create(dataInstance)

	// 因为gorm使用结构体，不更新零值，需要使用map更新零值
	reflectId := reflect.
		ValueOf(dataInstance).
		Elem().
		FieldByName("Id")
	if !reflectId.IsValid() {
		return ctx.JSON(200, message.Error("参数错误"))
	}

	id := int(reflectId.Int())
	db.Client.
		Model(&modelInstance).
		Where("id = ?", id).
		Updates(newData)

	// 保存后回调
	err = template.AfterSaved(ctx, id, data, model)

	return template.AfterSavedRedirectTo(ctx, err)
}
