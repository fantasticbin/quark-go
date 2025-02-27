// Code generated by hertz generator.

package main

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/adapter/hertzadapter"
	"github.com/quarkcloudio/quark-go/v3/app/admin"
	adminmodule "github.com/quarkcloudio/quark-go/v3/template/admin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	h := server.Default(server.WithHostPorts(":3000"))

	// 注册路由
	register(h)

	// 静态文件
	h.StaticFile("/admin/", "./web/app/admin/index.html")

	// WEB根目录
	fs := &app.FS{Root: "./web/app", IndexNames: []string{"index.html"}}
	h.StaticFS("/", fs)

	// 数据库配置信息
	dsn := "root:fK7xPGJi1gJfIief@tcp(127.0.0.1:3306)/quarkgo?charset=utf8&parseTime=True&loc=Local"

	// 配置资源
	config := &quark.Config{
		AppKey:    "123456",
		Providers: admin.Providers,
		DBConfig: &quark.DBConfig{
			Dialector: mysql.Open(dsn),
			Opts:      &gorm.Config{},
		},
	}

	// 创建对象
	b := quark.New(config)

	// 初始化安装
	adminmodule.Install()

	// 中间件
	b.Use(adminmodule.Middleware)

	// 适配hertz
	hertzadapter.Adapter(b, h)

	// 启动服务
	h.Spin()
}
