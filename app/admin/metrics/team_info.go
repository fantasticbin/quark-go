package metrics

import (
	"github.com/quarkcloudio/quark-go/v3/template/admin/component/descriptions"
	"github.com/quarkcloudio/quark-go/v3/template/admin/dashboard/metrics"
)

type TeamInfo struct {
	metrics.Descriptions
}

// 初始化
func (p *TeamInfo) Init() *TeamInfo {
	p.Title = "团队信息"
	p.Col = 12

	return p
}

// 计算数值
func (p *TeamInfo) Calculate() *descriptions.Component {

	field := &descriptions.Field{}

	return p.Init().Result([]interface{}{
		field.Text("作者").SetValue("tangtanglove"),
		field.Text("联系方式").SetValue("dai_hang_love@126.com"),
		field.Text("官方网址").SetValue("<a href='https://quarkcloud.io' target='_blank'>quarkcloud.io</a>"),
		field.Text("文档地址").SetValue("<a href='https://quarkcloud.io' target='_blank'>查看文档</a>"),
		field.Text("BUG反馈").SetValue("<a href='https://github.com/quarkcloudio/quark-go/v3/issues' target='_blank'>提交BUG</a>"),
		field.Text("代码仓储").SetValue("<a href='https://github.com/quarkcloudio/quark-go' target='_blank'>Github</a>"),
	})
}
