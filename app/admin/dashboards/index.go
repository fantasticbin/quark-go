package dashboards

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/app/admin/metrics"
	"github.com/quarkcloudio/quark-go/v3/template/admin/dashboard"
)

type Index struct {
	dashboard.Template
}

// 初始化
func (p *Index) Init(ctx *quark.Context) interface{} {
	p.Title = "仪表盘"

	return p
}

// 内容
func (p *Index) Cards(ctx *quark.Context) []interface{} {
	return []interface{}{
		&metrics.TotalAdmin{},
		&metrics.TotalLog{},
		&metrics.TotalPicture{},
		&metrics.TotalFile{},
		&metrics.SystemInfo{},
		&metrics.TeamInfo{},
	}
}
