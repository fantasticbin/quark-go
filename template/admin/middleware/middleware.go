package middleware

import (
	"strconv"
	"strings"

	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/app/admin/logins"
	"github.com/quarkcloudio/quark-go/v3/model"
)

// 中间件
func Handle(ctx *quark.Context) error {

	// 获取登录实例
	loginInstance := &logins.Index{}

	// 初始化路由
	loginInstance.RouteInit()

	// 加载自定义路由
	loginInstance.Route()

	// 获取登录模板定义的路由
	loginIndexRoutes := loginInstance.GetRouteMapping()

	inLoginRoute := false
	for _, v := range loginIndexRoutes {
		if v.Path == ctx.FullPath() {
			inLoginRoute = true
		}
	}

	// 排除登录路由
	if inLoginRoute {
		return ctx.Next()
	}

	// 排除非后台路由
	if !strings.Contains(ctx.Path(), "api/admin") {
		return ctx.Next()
	}

	// 定义管理员结构体
	adminInfo := &model.UserClaims{}

	// 获取登录管理员信息
	err := ctx.JwtAuthUser(adminInfo)
	if err != nil {
		return ctx.JSON(401, quark.Error(err.Error()))
	}

	guardName := adminInfo.GuardName
	if guardName != "admin" {
		return ctx.JSON(401, quark.Error("401 Unauthozied"))
	}

	// 管理员id
	if adminInfo.Id != 1 {
		result1, err := (&model.CasbinRule{}).Enforce("admin|"+strconv.Itoa(adminInfo.Id), ctx.FullPath(), "Any")
		if err != nil {
			return ctx.JSON(500, quark.Error(err.Error()))
		}

		result2, err := (&model.CasbinRule{}).Enforce("admin|"+strconv.Itoa(adminInfo.Id), ctx.FullPath(), ctx.Method())
		if err != nil {
			return ctx.JSON(500, quark.Error(err.Error()))
		}

		result3, err := (&model.CasbinRule{}).Enforce("admin|"+strconv.Itoa(adminInfo.Id), ctx.Path(), "Any")
		if err != nil {
			return ctx.JSON(500, quark.Error(err.Error()))
		}

		result4, err := (&model.CasbinRule{}).Enforce("admin|"+strconv.Itoa(adminInfo.Id), ctx.Path(), ctx.Method())
		if err != nil {
			return ctx.JSON(500, quark.Error(err.Error()))
		}

		if !(result1 || result2 || result3 || result4) {
			return ctx.JSON(403, quark.Error("403 Forbidden"))
		}
	}

	// 记录操作日志
	(&model.ActionLog{}).InsertGetId(model.ActionLog{
		ObjectId: adminInfo.Id,
		Url:      ctx.Path(),
		Ip:       ctx.ClientIP(),
		Type:     "admin",
	})

	return ctx.Next()
}
