package requests

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/message"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
	"gorm.io/gorm"
)

type ActionRequest struct{}

// 执行行为
func (p *ActionRequest) Handle(ctx *quark.Context) error {
	var result error

	// 模版实例
	template := ctx.Template.(types.Resourcer)

	// 模型结构体
	modelInstance := template.GetModel()

	// Gorm对象
	model := db.Client.Model(modelInstance)

	// 查询条件
	model = template.BuildActionQuery(ctx, model)

	actions := template.Actions(ctx)
	for _, v := range actions {
		actionInstance := v.(types.Actioner)

		// 初始化模版
		actionInstance.TemplateInit(ctx)

		// 初始化
		actionInstance.Init(ctx)

		// uri唯一标识
		uriKey := actionInstance.GetUriKey(v)

		// 获取行为类型
		actionType := actionInstance.GetActionType()

		if actionType == "dropdown" {
			dropdownActioner := v.(types.Dropdowner)
			for _, dropdownAction := range dropdownActioner.GetActions() {
				uriKey := dropdownActioner.GetUriKey(dropdownAction)
				if ctx.Param("uriKey") == uriKey {
					result = dropdownAction.(interface {
						Handle(*quark.Context, *gorm.DB) error
					}).Handle(ctx, model)

					// 执行完后回调
					err := template.AfterAction(ctx, uriKey, model)
					if err != nil {
						return err
					}

					return result
				}
			}
		} else {
			if ctx.Param("uriKey") == uriKey {
				result = v.(interface {
					Handle(*quark.Context, *gorm.DB) error
				}).Handle(ctx, model)

				// 执行完后回调
				err := template.AfterAction(ctx, uriKey, model)
				if err != nil {
					return err
				}

				return result
			}
		}
	}

	return result
}

// 行为表单值
func (p *ActionRequest) Values(ctx *quark.Context) error {
	var data map[string]interface{}

	// 模版实例
	template := ctx.Template.(types.Resourcer)

	// 解析行为
	actions := template.Actions(ctx)
	for _, v := range actions {

		actionInstance := v.(types.Actioner)

		// 初始化模版
		actionInstance.TemplateInit(ctx)

		// 初始化
		actionInstance.Init(ctx)

		// uri唯一标识
		uriKey := actionInstance.GetUriKey(v)

		// 获取行为类型
		actionType := actionInstance.GetActionType()
		if actionType == "dropdown" {
			dropdownActioner := v.(types.Dropdowner)
			for _, dropdownAction := range dropdownActioner.GetActions() {
				uriKey := dropdownActioner.GetUriKey(dropdownAction)
				if ctx.Param("uriKey") == uriKey {
					data = dropdownAction.(interface {
						Data(*quark.Context) map[string]interface{}
					}).Data(ctx)
				}
			}
		} else {
			if ctx.Param("uriKey") == uriKey {
				data = v.(interface {
					Data(*quark.Context) map[string]interface{}
				}).Data(ctx)
			}
		}
	}

	return ctx.JSON(200, message.Success("获取成功", "", data))
}
