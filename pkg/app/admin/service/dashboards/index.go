package dashboards

import (
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/service/metrics"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/template/dashboard"
	"github.com/quarkcloudio/quark-go/v3/pkg/builder"
)

type Index struct {
	dashboard.Template
}

// 初始化
func (p *Index) Init(ctx *builder.Context) interface{} {
	p.Title = "仪表盘"

	return p
}

// 内容
func (p *Index) Cards(ctx *builder.Context) []interface{} {
	return []interface{}{
		&metrics.TotalAdmin{},
		&metrics.TotalLog{},
		&metrics.TotalPicture{},
		&metrics.TotalFile{},
		&metrics.SystemInfo{},
		&metrics.TeamInfo{},
	}
}
