package middleware

import (
	"strings"

	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/model"
)

// 中间件
func Handle(ctx *quark.Context) error {

	// 排除非后台路由
	if !strings.Contains(ctx.Path(), "api/miniapp/user") {
		return ctx.Next()
	}

	// 获取登录信息
	userInfo, err := (&model.User{}).GetAuthUser(ctx.Engine.GetConfig().AppKey, ctx.Token())
	if err != nil {
		return ctx.JSON(401, quark.Error(err.Error()))
	}

	guardName := userInfo.GuardName
	if guardName != "user" {
		return ctx.JSON(401, quark.Error("401 Unauthozied"))
	}

	return ctx.Next()
}
