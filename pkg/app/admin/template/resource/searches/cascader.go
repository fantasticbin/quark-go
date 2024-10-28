package searches

import (
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/component/form/fields/cascader"
	"github.com/quarkcloudio/quark-go/v3/pkg/builder"
)

type Cascader struct {
	Search
	CascaderOptions []cascader.Option
}

// 初始化模板
func (p *Cascader) TemplateInit(ctx *builder.Context) interface{} {
	p.Component = "cascaderField"

	return p
}

// 可选项数据源
//
// SetOptions([]cascader.Option {{Value :"zhejiang", Label:"Zhejiang"}})
//
// 或者
//
// SetOptions(options, "parent_key_name", "label_name", "value_name")
//
// 或者
//
// SetOptions(options, 0, "parent_key_name", "label_name", "value_name")
func (p *Cascader) SetOptions(options ...interface{}) *Cascader {
	if len(options) == 1 {
		getOptions, ok := options[0].([]cascader.Option)
		if ok {
			p.CascaderOptions = getOptions
			return p
		}
	}
	if len(options) == 4 {
		p.CascaderOptions = cascader.New().ListToOptions(options[0], 0, options[1].(string), options[2].(string), options[3].(string))
	}
	if len(options) == 5 {
		p.CascaderOptions = cascader.New().ListToOptions(options[0], options[1].(int), options[2].(string), options[3].(string), options[4].(string))
	}
	return p
}

// 设置Option
func (p *Cascader) Option(label string, value interface{}) cascader.Option {

	return cascader.Option{
		Value: value,
		Label: label,
	}
}
