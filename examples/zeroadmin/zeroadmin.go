package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/adapter/zeroadapter"
	"github.com/quarkcloudio/quark-go/v3/app/admin"
	"github.com/quarkcloudio/quark-go/v3/examples/zeroadmin/internal/config"
	"github.com/quarkcloudio/quark-go/v3/examples/zeroadmin/internal/handler"
	"github.com/quarkcloudio/quark-go/v3/examples/zeroadmin/internal/svc"
	adminmodule "github.com/quarkcloudio/quark-go/v3/template/admin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/zeroadmin-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// WEB根目录
	staticFile("/", "./web/app", server)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

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

	// 适配gozero
	zeroadapter.Adapter(b, server)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// 加载静态文件
func staticFile(root string, dirPath string, server *rest.Server) {
	rd, _ := os.ReadDir(dirPath)

	for _, f := range rd {
		fileName := f.Name()
		subPath := root + fileName + "/"
		subDirPath := dirPath + "/" + fileName
		if isDir(subDirPath) {
			staticFile(subPath, subDirPath, server)
		}
	}

	server.AddRoute(
		rest.Route{
			Method:  http.MethodGet,
			Path:    root + ":file",
			Handler: http.StripPrefix(root, http.FileServer(http.Dir(dirPath))).ServeHTTP,
		},
	)
}

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
