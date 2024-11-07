package searches

import (
	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/searches"
	"gorm.io/gorm"
)

type DateTimeRangeField struct {
	searches.DatetimeRange
}

// 日期时间范围
func DatetimeRange(column string, name string) *DateTimeRangeField {
	field := &DateTimeRangeField{}

	field.Column = column
	field.Name = name

	return field
}

// 执行查询
func (p *DateTimeRangeField) Apply(ctx *quark.Context, query *gorm.DB, value interface{}) *gorm.DB {
	values, ok := value.([]interface{})
	if !ok {
		return query
	}

	return query.Where(p.Column+" BETWEEN ? AND ?", values[0], values[1])
}
