package actions

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/message"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/actions"
	"gorm.io/gorm"
)

type SelectOptionsAction struct {
	actions.Action
}

// 执行行为句柄
func (p *SelectOptionsAction) Handle(ctx *quark.Context, query *gorm.DB) error {
	return ctx.JSON(200, message.Success("操作成功"))
}
