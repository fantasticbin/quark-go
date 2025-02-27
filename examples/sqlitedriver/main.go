package main

import (
	"github.com/quarkcloudio/quark-go/v3"
	adminservice "github.com/quarkcloudio/quark-go/v3/app/admin"
	toolservice "github.com/quarkcloudio/quark-go/v3/app/tool"
	adminmodule "github.com/quarkcloudio/quark-go/v3/template/admin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	// 定义服务
	var providers []interface{}

	// 数据库配置信息
	dsn := "./data.db"

	// 加载后台服务
	providers = append(providers, adminservice.Providers...)

	// 加载工具服务
	providers = append(providers, toolservice.Providers...)

	// 配置资源
	config := &quark.Config{
		AppKey:    "123456",
		Providers: providers,
		DBConfig: &quark.DBConfig{
			Dialector: sqlite.Open(dsn),
			Opts:      &gorm.Config{},
		},
	}

	// 实例化对象
	b := quark.New(config)

	// WEB根目录
	b.Static("/", "./web/app")

	// 初始化安装
	adminmodule.Install()

	// 中间件
	b.Use(adminmodule.Middleware)

	// 响应Get请求
	b.GET("/", func(ctx *quark.Context) error {
		return ctx.String(200, "Hello World!")
	})

	// 启动服务
	b.Run(":3000")
}
