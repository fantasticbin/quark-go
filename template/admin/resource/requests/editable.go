package requests

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/message"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
)

type EditableRequest struct{}

// 执行行为
func (p *EditableRequest) Handle(ctx *quark.Context) error {
	var (
		id    interface{}
		field string
		value interface{}
	)

	// 获取所有Query数据
	data := ctx.AllQuerys()
	if data == nil {
		return ctx.JSON(200, message.Error("参数错误！"))
	}

	id = data["id"]
	if id == nil {
		return ctx.JSON(200, message.Error("id不能为空！"))
	}

	// 模版实例
	template := ctx.Template.(types.Resourcer)

	// 获取模型结构体
	modelInstance := template.GetModel()

	// 创建Gorm对象
	model := db.Client.Model(&modelInstance)

	// 解析数据
	for k, v := range data {
		if v == "true" {
			v = 1
		} else if v == "false" {
			v = 0
		}

		if k != "id" && k != "_t" {
			field = k
			value = v
		}
	}

	if field == "" {
		return ctx.JSON(200, message.Error("参数错误！"))
	}

	if value == nil {
		return ctx.JSON(200, message.Error("参数错误！"))
	}

	// 创建表格行内编辑查询
	query := template.BuildEditableQuery(ctx, model)

	// 更新数据
	err := query.Update(field, value).Error
	if err != nil {
		return ctx.JSON(200, message.Error(err.Error()))
	}

	// 行为执行后回调
	result := template.AfterEditable(ctx, id, field, value)
	if result != nil {
		return result
	}

	return ctx.JSON(200, message.Success("操作成功"))
}
