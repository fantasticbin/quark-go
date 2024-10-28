package searches

import (
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
