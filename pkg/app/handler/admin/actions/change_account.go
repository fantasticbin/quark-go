package actions

import (
	"encoding/json"

	"github.com/quarkcms/quark-go/pkg/app/model"
	"github.com/quarkcms/quark-go/pkg/builder"
	"github.com/quarkcms/quark-go/pkg/builder/template/adminresource/actions"
	"github.com/quarkcms/quark-go/pkg/hash"
	"gorm.io/gorm"
)

type ChangeAccount struct {
	actions.Action
}

// 执行行为句柄
func (p *ChangeAccount) Handle(ctx *builder.Context, query *gorm.DB) error {
	data := map[string]interface{}{}
	json.Unmarshal(ctx.Body(), &data)
	if data["avatar"] != "" {
		data["avatar"], _ = json.Marshal(data["avatar"])
	} else {
		data["avatar"] = nil
	}

	// 加密密码
	if data["password"] != nil {
		data["password"] = hash.Make(data["password"].(string))
	}

	// 获取登录管理员信息
	adminInfo, err := (&model.Admin{}).GetAuthUser(ctx.Engine.GetConfig().AppKey, ctx.Token())
	if err != nil {
		return ctx.JSONError(err.Error())
	}

	err = query.Where("id", adminInfo.Id).Updates(data).Error
	if err != nil {
		return ctx.JSONError(err.Error())
	}

	return ctx.JSONOk("操作成功")
}
