package searches

import (
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/component/form/fields/radio"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/template/resource/searches"
	"github.com/quarkcloudio/quark-go/v3/pkg/builder"
	"gorm.io/gorm"
)

type RadioField struct {
	searches.Radio
}

// 下拉框
func Radio(column string, name string) *RadioField {
	field := &RadioField{}

	field.Column = column
	field.Name = name

	return field
}

// 执行查询
func (p *RadioField) Apply(ctx *builder.Context, query *gorm.DB, value interface{}) *gorm.DB {
	return query.Where(p.Column+" = ?", value)
}

// 属性
func (p *RadioField) Options(ctx *builder.Context) interface{} {
	return p.RadioOptions
}

// 设置属性，示例：[]radio.Option{{Value: 1, Label: "男"}, {Value: 2, Label: "女"}}
//
// 或者
//
// SetOptions(options, "label_name", "value_name")
func (p *RadioField) SetOptions(options ...interface{}) *RadioField {
	if len(options) == 1 {
		getOptions, ok := options[0].([]radio.Option)
		if ok {
			p.RadioOptions = getOptions
			return p
		}
	}
	if len(options) == 3 {
		p.RadioOptions = radio.New().ListToOptions(options[0], options[1].(string), options[2].(string))
	}
	return p
}
