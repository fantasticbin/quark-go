package model

import (
	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/utils/datetime"
)

// 字段
type Config struct {
	Id        int               `json:"id" gorm:"autoIncrement"`
	Title     string            `json:"title" gorm:"size:255;not null"`
	Type      string            `json:"type" gorm:"size:20;not null"`
	Name      string            `json:"name" gorm:"size:255;not null"`
	Sort      int               `json:"sort" gorm:"size:11;default:0"`
	GroupName string            `json:"group_name" gorm:"size:255;not null"`
	Value     string            `json:"value" gorm:"size:2000"`
	Remark    string            `json:"remark" gorm:"size:100;not null"`
	Status    int               `json:"status" gorm:"size:1;not null;default:1"`
	CreatedAt datetime.Datetime `json:"created_at"`
	UpdatedAt datetime.Datetime `json:"updated_at"`
}

// 配置表
func (model *Config) Seeder() {
	seeders := []Config{
		{Title: "网站名称", Type: "text", Name: "WEB_SITE_NAME", Sort: 0, GroupName: "基本", Value: "QuarkCloud", Remark: "", Status: 1},
		{Title: "关键字", Type: "text", Name: "WEB_SITE_KEYWORDS", Sort: 0, GroupName: "基本", Value: "QuarkCloud", Remark: "", Status: 1},
		{Title: "描述", Type: "textarea", Name: "WEB_SITE_DESCRIPTION", Sort: 0, GroupName: "基本", Value: "QuarkCloud", Remark: "", Status: 1},
		{Title: "Logo", Type: "picture", Name: "WEB_SITE_LOGO", Sort: 0, GroupName: "基本", Value: "", Remark: "", Status: 1},
		{Title: "统计代码", Type: "textarea", Name: "WEB_SITE_SCRIPT", Sort: 0, GroupName: "基本", Value: "", Remark: "", Status: 1},
		{Title: "网站域名", Type: "text", Name: "WEB_SITE_DOMAIN", Sort: 0, GroupName: "基本", Value: "", Remark: "", Status: 1},
		{Title: "网站版权", Type: "text", Name: "WEB_SITE_COPYRIGHT", Sort: 0, GroupName: "基本", Value: "© Company 2018", Remark: "", Status: 1},
		{Title: "开启SSL", Type: "switch", Name: "SSL_OPEN", Sort: 0, GroupName: "基本", Value: "0", Remark: "", Status: 1},
		{Title: "开启网站", Type: "switch", Name: "WEB_SITE_OPEN", Sort: 0, GroupName: "基本", Value: "1", Remark: "", Status: 1},
		{Title: "KeyID", Type: "text", Name: "OSS_ACCESS_KEY_ID", Sort: 0, GroupName: "阿里云存储", Value: "", Remark: "你的AccessKeyID", Status: 1},
		{Title: "KeySecret", Type: "text", Name: "OSS_ACCESS_KEY_SECRET", Sort: 0, GroupName: "阿里云存储", Value: "", Remark: "你的AccessKeySecret", Status: 1},
		{Title: "EndPoint", Type: "text", Name: "OSS_ENDPOINT", Sort: 0, GroupName: "阿里云存储", Value: "", Remark: "地域节点", Status: 1},
		{Title: "Bucket域名", Type: "text", Name: "OSS_BUCKET", Sort: 0, GroupName: "阿里云存储", Value: "", Remark: "", Status: 1},
		{Title: "自定义域名", Type: "text", Name: "OSS_MYDOMAIN", Sort: 0, GroupName: "阿里云存储", Value: "", Remark: "例如：oss.web.com", Status: 1},
		{Title: "开启云存储", Type: "switch", Name: "OSS_OPEN", Sort: 0, GroupName: "阿里云存储", Value: "0", Remark: "", Status: 1},
	}

	db.Client.Create(&seeders)
}
