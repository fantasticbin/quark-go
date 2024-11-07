package resource

import (
	"encoding/json"
	"strings"

	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/template/admin/resource/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 创建行为查询
func (p *Template) BuildActionQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	template := ctx.Template.(types.Resourcer)

	// 初始化查询
	query = p.initializeQuery(ctx, query)

	// 执行查询，这里使用的是透传的实例
	query = template.ActionQuery(ctx, query)

	return query
}

// 创建详情页查询
func (p *Template) BuildDetailQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	template := ctx.Template.(types.Resourcer)

	// 初始化查询
	query = p.initializeQuery(ctx, query)

	// 执行列表查询，这里使用的是透传的实例
	query = template.DetailQuery(ctx, query)

	return query
}

// 创建编辑页查询
func (p *Template) BuildEditQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	template := ctx.Template.(types.Resourcer)

	// 初始化查询
	query = p.initializeQuery(ctx, query)

	// 执行查询，这里使用的是透传的实例
	query = template.EditQuery(ctx, query)

	return query
}

// 创建表格行内编辑查询
func (p *Template) BuildEditableQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	template := ctx.Template.(types.Resourcer)

	// 初始化查询
	query = p.initializeQuery(ctx, query)

	// 执行查询，这里使用的是透传的实例
	query = template.EditableQuery(ctx, query)

	return query
}

// 创建导出查询
func (p *Template) BuildExportQuery(ctx *quark.Context, query *gorm.DB, search []interface{}, filters []interface{}, columnFilters map[string]interface{}, orderings map[string]interface{}) *gorm.DB {
	template := ctx.Template.(types.Resourcer)

	// 初始化查询
	query = p.initializeQuery(ctx, query)

	// 执行列表查询，这里使用的是透传的实例
	query = template.ExportQuery(ctx, query)

	// 执行搜索查询
	query = p.applySearch(ctx, query, search)

	// 执行过滤器查询
	query = p.applyFilters(query, filters)

	// 执行表格列上过滤器查询
	query = p.applyColumnFilters(query, columnFilters)

	// 获取排序规则
	defaultOrder := template.GetQueryOrder()
	if defaultOrder == "" {
		defaultOrder = template.GetExportQueryOrder()
	}

	if defaultOrder == "" {
		defaultOrder = "id desc"
	}

	// 执行排序查询
	query = p.applyOrderings(query, orderings, defaultOrder)

	return query
}

// 创建列表查询
func (p *Template) BuildIndexQuery(ctx *quark.Context, query *gorm.DB, search []interface{}, filters []interface{}, columnFilters map[string]interface{}, orderings map[string]interface{}) *gorm.DB {
	template := ctx.Template.(types.Resourcer)

	// 初始化查询
	query = p.initializeQuery(ctx, query)

	// 执行列表查询，这里使用的是透传的实例
	query = template.IndexQuery(ctx, query)

	// 执行搜索查询
	query = p.applySearch(ctx, query, search)

	// 执行过滤器查询
	query = p.applyFilters(query, filters)

	// 执行表格列上过滤器查询
	query = p.applyColumnFilters(query, columnFilters)

	// 获取排序规则
	defaultOrder := template.GetQueryOrder()
	if defaultOrder == "" {
		defaultOrder = template.GetIndexQueryOrder()
	}

	if defaultOrder == "" {
		defaultOrder = "id desc"
	}

	// 执行排序查询
	query = p.applyOrderings(query, orderings, defaultOrder)

	return query
}

// 创建更新查询
func (p *Template) BuildUpdateQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	template := ctx.Template.(types.Resourcer)

	// 初始化查询
	query = p.initializeQuery(ctx, query)

	// 执行查询，这里使用的是透传的实例
	query = template.UpdateQuery(ctx, query)

	return query
}

// 初始化查询
func (p *Template) initializeQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	template := ctx.Template.(types.Resourcer)

	return template.Query(ctx, query)
}

// 执行搜索表单查询
func (p *Template) applySearch(ctx *quark.Context, query *gorm.DB, search []interface{}) *gorm.DB {
	querys := ctx.AllQuerys()
	var data map[string]interface{}
	if querys["search"] == nil {
		return query
	}
	err := json.Unmarshal([]byte(querys["search"].(string)), &data)
	if err != nil {
		return query
	}
	for _, v := range search {
		// 获取字段
		column := v.(interface {
			GetColumn(search interface{}) string
		}).GetColumn(v) // 字段名，支持数组
		value := data[column]
		if value != nil {
			query = v.(interface {
				Apply(*quark.Context, *gorm.DB, interface{}) *gorm.DB
			}).Apply(ctx, query, value)
		}
	}

	return query
}

// 执行表格列上过滤器查询
func (p *Template) applyColumnFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if len(filters) == 0 || filters == nil {
		return query
	}
	for k, v := range filters {
		if v != nil {
			query = query.Where(k+" IN ?", v)
		}
	}

	return query
}

// 执行过滤器查询
func (p *Template) applyFilters(query *gorm.DB, filters []interface{}) *gorm.DB {
	// todo
	return query
}

// 执行排序查询
func (p *Template) applyOrderings(query *gorm.DB, orderings map[string]interface{}, defaultOrder string) *gorm.DB {
	if len(orderings) == 0 || orderings == nil {
		return query.Order(defaultOrder)
	}
	var order clause.OrderByColumn
	for key, v := range orderings {
		if v != nil {
			if v == "descend" {
				order = clause.OrderByColumn{Column: clause.Column{Name: key}, Desc: true}
			} else {
				order = clause.OrderByColumn{Column: clause.Column{Name: key}, Desc: false}
			}
			query = query.Order(order)
		}
	}

	return query
}

// 全局查询
func (p *Template) Query(ctx *quark.Context, query *gorm.DB) *gorm.DB {

	return query
}

// 行为查询
func (p *Template) ActionQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	id := ctx.Query("id", "")
	if id != "" {
		if strings.Contains(id.(string), ",") {
			query.Where("id IN ?", strings.Split(id.(string), ","))
		} else {
			query.Where("id = ?", id)
		}
	}

	return query
}

// 详情查询
func (p *Template) DetailQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	id := ctx.Query("id", "")
	if id != "" {
		query.Where("id = ?", id)
	}

	return query
}

// 编辑查询
func (p *Template) EditQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	id := ctx.Query("id", "")
	if id != "" {
		query.Where("id = ?", id)
	}

	return query
}

// 表格行内编辑查询
func (p *Template) EditableQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	data := ctx.AllQuerys()
	if data != nil {
		if data["id"] != nil {
			query.Where("id = ?", data["id"])
		}
	}

	return query
}

// 导出查询
func (p *Template) ExportQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {

	return query
}

// 列表查询
func (p *Template) IndexQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {

	return query
}

// 更新查询
func (p *Template) UpdateQuery(ctx *quark.Context, query *gorm.DB) *gorm.DB {
	data := map[string]interface{}{}
	ctx.Bind(&data)
	if data != nil {
		if data["id"] != nil {
			query.Where("id = ?", data["id"])
		}
	}

	return query
}
