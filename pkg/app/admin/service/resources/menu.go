package resources

import (
	"strconv"

	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/component/form/fields/radio"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/component/form/rule"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/model"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/service/actions"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/service/searches"
	"github.com/quarkcloudio/quark-go/v3/pkg/app/admin/template/resource"
	"github.com/quarkcloudio/quark-go/v3/pkg/builder"
	"github.com/quarkcloudio/quark-go/v3/pkg/utils/lister"
	"gorm.io/gorm"
)

type Menu struct {
	resource.Template
}

// 初始化
func (p *Menu) Init(ctx *builder.Context) interface{} {

	// 标题
	p.Title = "菜单"

	// 模型
	p.Model = &model.Menu{}

	// 分页
	p.PerPage = false

	// 默认排序
	p.QueryOrder = "sort asc"

	return p
}

// 字段
func (p *Menu) Fields(ctx *builder.Context) []interface{} {
	field := &resource.Field{}

	// 权限列表
	permissions, _ := (&model.Permission{}).DataSource()

	// 菜单列表
	menus, _ := (&model.Menu{}).GetListWithRoot()

	return []interface{}{
		field.Hidden("id", "ID"),                 // 列表读取且不展示的字段
		field.Hidden("pid", "PID").OnlyOnIndex(), // 列表读取且不展示的字段
		field.Group([]interface{}{
			field.Text("name", "名称").
				SetRules([]rule.Rule{
					rule.Required("名称必须填写"),
				}),
			field.Text("guard_name", "守卫").
				SetRules([]rule.Rule{
					rule.Required("守卫必须填写"),
				}).
				SetDefault("admin").
				OnlyOnForms(),
			field.Icon("icon", "图标").OnlyOnForms(),
		}),
		field.Group([]interface{}{
			field.Number("sort", "排序").
				SetEditable(true).
				SetDefault(0),
			field.TreeSelect("pid", "父节点").
				SetTreeData(menus, -1, "pid", "name", "id").
				SetDefault(0).
				OnlyOnForms(),
			field.Switch("status", "状态").
				SetTrueValue("正常").
				SetFalseValue("禁用").
				SetEditable(true).
				SetDefault(true),
		}),
		field.Group([]interface{}{
			field.Radio("type", "类型").
				SetOptions([]radio.Option{
					field.RadioOption("目录", 1),
					field.RadioOption("菜单", 2),
					field.RadioOption("按钮", 3),
				}).
				SetRules([]rule.Rule{
					rule.Required("类型必须选择"),
				}).
				SetDefault(1),
			field.Switch("show", "显示").
				SetTrueValue("显示").
				SetFalseValue("隐藏").
				SetEditable(true).
				SetDefault(true),
		}),
		field.Dependency().
			SetWhen("type", 1, func() interface{} {
				return []interface{}{
					field.Text("path", "路由").
						SetRules([]rule.Rule{
							rule.Required("路由必须填写"),
						}).
						SetEditable(true).
						SetHelp("前端路由").
						SetWidth(400),
				}
			}),
		field.Dependency().
			SetWhen("type", 2, func() interface{} {
				return []interface{}{
					field.Switch("is_engine", "引擎组件").
						SetTrueValue("是").
						SetFalseValue("否").
						SetDefault(true),
					field.Switch("is_link", "外部链接").
						SetTrueValue("是").
						SetFalseValue("否").
						SetDefault(false),
					field.Text("path", "路由").
						SetRules([]rule.Rule{
							rule.Required("路由必须填写"),
						}).
						SetEditable(true).
						SetHelp("前端路由或后端api").
						SetWidth(400).
						OnlyOnForms(),
				}
			}),
		field.Dependency().
			SetWhen("type", 3, func() interface{} {
				return []interface{}{
					field.Transfer("permission_ids", "绑定权限").
						SetDataSource(permissions).
						SetListStyle(map[string]interface{}{
							"width":  320,
							"height": 300,
						}).
						SetShowSearch(true).
						OnlyOnForms(),
				}
			}),
	}
}

// 搜索
func (p *Menu) Searches(ctx *builder.Context) []interface{} {
	return []interface{}{
		searches.Input("name", "名称"),
		searches.Input("path", "路由"),
		searches.Status(),
	}
}

// 行为
func (p *Menu) Actions(ctx *builder.Context) []interface{} {
	return []interface{}{
		actions.MenuCreateDrawer(),
		actions.BatchDelete(),
		actions.BatchDisable(),
		actions.BatchEnable(),
		actions.ChangeStatus(),
		actions.MenuEditDrawer(),
		actions.Delete(),
		actions.FormSubmit(),
		actions.FormReset(),
		actions.FormBack(),
		actions.FormExtraBack(),
	}
}

// 列表页面显示前回调
func (p *Menu) BeforeIndexShowing(ctx *builder.Context, list []map[string]interface{}) []interface{} {
	data := ctx.AllQuerys()
	if search, ok := data["search"].(map[string]interface{}); ok && search != nil {
		result := []interface{}{}
		for _, v := range list {
			result = append(result, v)
		}
		return result
	}
	// 转换成树形表格
	tree, _ := lister.ListToTree(list, "id", "pid", "children", 0)
	return tree
}

// 编辑页面显示前回调
func (p *Menu) BeforeEditing(ctx *builder.Context, data map[string]interface{}) map[string]interface{} {
	id := ctx.Query("id", "")
	idInt, err := strconv.Atoi(id.(string))
	if id != "" && err == nil {
		permissionIds := []int{}
		permissions, err := (&model.CasbinRule{}).GetMenuPermissions(idInt)
		if err == nil {
			for _, v := range permissions {
				permissionIds = append(permissionIds, v.Id)
			}
		}
		data["permission_ids"] = permissionIds
	}
	return data
}

// 保存后回调
func (p *Menu) AfterSaved(ctx *builder.Context, id int, data map[string]interface{}, result *gorm.DB) error {
	if data["permission_ids"] != nil {
		err := (&model.CasbinRule{}).AddMenuPermission(id, data["permission_ids"])
		if err != nil {
			return err
		}
	}
	return result.Error
}
