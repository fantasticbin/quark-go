package searches

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/searches"
	"gorm.io/gorm"
)

type DatetimeField struct {
	searches.DatetimeRange
}

// 日期时间
func Datetime(column string, name string) *DatetimeField {
	field := &DatetimeField{}

	field.Column = column
	field.Name = name

	return field
}

// 执行查询
func (p *DatetimeField) Apply(ctx *quark.Context, query *gorm.DB, value interface{}) *gorm.DB {
	return query.Where(p.Column+" = ?", value)
}
