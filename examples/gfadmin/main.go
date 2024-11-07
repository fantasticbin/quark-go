package main

import (
	"os"
	"path"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/glebarez/sqlite"
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/adapter/gfadapter"
	"github.com/quarkcloudio/quark-go/v3/app/admin"
	adminmodule "github.com/quarkcloudio/quark-go/v3/template/admin"
	"gorm.io/gorm"
)

func main() {
	currentDir, _ := os.Getwd()
	s := g.Server()

	// 配置资源
	config := &quark.Config{

		// JWT加密密串
		AppKey: "123456",

		// 加载服务
		Providers: admin.Providers,

		// 数据库配置
		DBConfig: &quark.DBConfig{
			Dialector: sqlite.Open("./examples/gfadmin/data.db"),
			Opts:      &gorm.Config{},
		},
	}

	// 创建对象
	b := quark.New(config)

	// 初始化安装
	adminmodule.Install()

	// 中间件
	b.Use(adminmodule.Middleware)

	// WEB根目录
	s.SetServerRoot("./web/app")
	s.AddSearchPath(path.Join(currentDir, "web/app/admin"))

	// 适配goframe
	gfadapter.Adapter(b, s)

	s.SetPort(3000)
	s.Run()
}
