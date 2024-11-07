package service

import (
	"github.com/quarkcloudio/quark-go/v3/app/miniapp/forms"
	"github.com/quarkcloudio/quark-go/v3/app/miniapp/layouts"
	"github.com/quarkcloudio/quark-go/v3/app/miniapp/pages"
)

// 注册服务
var Providers = []interface{}{
	&layouts.Index{},
	&pages.Index{},
	&pages.My{},
	&forms.Demo{},
}
