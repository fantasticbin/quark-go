package service

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/model"
)

type PictureService struct{}

// / 初始化
func NewPictureService() *PictureService {
	return &PictureService{}
}

// 获取列表
func (p *PictureService) GetListBySearch(appKey string, tokenString string, categoryId interface{}, name interface{}, startDate interface{}, endDate interface{}, page int) (list []model.Picture, total int64, Error error) {
	pictures := []model.Picture{}

	adminInfo, err := NewUserService().GetAuthUser(appKey, tokenString)
	if err != nil {
		return pictures, 0, err
	}

	query := db.Client.Model(&model.Picture{}).
		Where("status =?", 1).
		Where("obj_type = ?", "ADMIN").
		Where("obj_id", adminInfo.Id)

	if categoryId != "" {
		query.Where("picture_category_id =?", categoryId)
	}
	if name != "" {
		query.Where("name LIKE %?%", name)
	}
	if startDate != "" && endDate != "" {
		query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	query.Count(&total)
	query.
		Order("id desc").
		Limit(8).
		Offset((page - 1) * 8).
		Find(&pictures)

	for k, v := range pictures {
		v.Url = p.GetPath(v.Url) + "?timestamp=" + strconv.Itoa(int(time.Now().Unix()))
		pictures[k] = v
	}

	return pictures, total, nil
}

// 插入数据并返回ID
func (p *PictureService) InsertGetId(picture model.Picture) (id int, Error error) {
	err := db.Client.Create(&picture).Error
	if err != nil {
		return id, err
	}

	return picture.Id, nil
}

// 通过Id删除记录
func (p *PictureService) DeleteById(id interface{}) error {

	return db.Client.Model(model.Picture{}).Where("id =?", id).Delete("").Error
}

// 根据id查询文件信息
func (p *PictureService) GetInfoById(id interface{}) (picture model.Picture, Error error) {
	err := db.Client.Where("status = ?", 1).Where("id = ?", id).First(&picture).Error

	return picture, err
}

// 根据id更新文件信息
func (p *PictureService) UpdateById(id interface{}, data model.Picture) (Error error) {
	err := db.Client.Where("status = ?", 1).Where("id = ?", id).Updates(&data).Error

	return err
}

// 根据hash查询文件信息
func (p *PictureService) GetInfoByHash(hash string) (picture model.Picture, Error error) {
	err := db.Client.Where("status = ?", 1).Where("hash = ?", hash).First(&picture).Error

	return picture, err
}

// 获取图片路径
func (p *PictureService) GetPath(id interface{}) string {
	http, path := "", ""
	webSiteDomain := NewConfigService().GetValue("WEB_SITE_DOMAIN")
	WebConfig := NewConfigService().GetValue("SSL_OPEN")
	if webSiteDomain != "" {
		if WebConfig == "1" {
			http = "https://"
		} else {
			http = "http://"
		}
	}

	if getId, ok := id.(string); ok {
		if strings.Contains(getId, "//") && !strings.Contains(getId, "{") {
			return getId
		}
		if strings.Contains(getId, "./") && !strings.Contains(getId, "{") {
			return http + webSiteDomain + strings.Replace(getId, "./web/app/", "/", -1)
		}
		if strings.Contains(getId, "/") && !strings.Contains(getId, "{") {
			return http + webSiteDomain + getId
		}

		// json字符串
		if strings.Contains(getId, "{") {
			var jsonData interface{}
			json.Unmarshal([]byte(getId), &jsonData)
			if mapData, ok := jsonData.(map[string]interface{}); ok {
				path = mapData["url"].(string)
			}

			// 如果为数组，返回第一个key的path
			if arrayData, ok := jsonData.([]map[string]interface{}); ok {
				path = arrayData[0]["url"].(string)
			}
		}
		if strings.Contains(path, "//") {
			return path
		}
		if strings.Contains(path, "./") {
			path = strings.Replace(path, "./web/app/", "/", -1)
		}
		if path != "" {
			// 如果设置域名，则加上域名前缀
			return http + webSiteDomain + path
		}
	}

	picture := model.Picture{}
	db.Client.Where("id", id).Where("status", 1).First(&picture)
	if picture.Id != 0 {
		path = picture.Url
		if strings.Contains(path, "//") {
			return path
		}
		if strings.Contains(path, "./") {
			path = strings.Replace(path, "./web/app/", "/", -1)
		}
	}
	if path != "" {
		// 如果设置域名，则加上域名前缀
		return http + webSiteDomain + path
	}

	return http + webSiteDomain + "/admin/default.png"
}

// 获取多图片路径
func (p *PictureService) GetPaths(id interface{}) []string {
	var paths []string
	http, path := "", ""
	webSiteDomain := NewConfigService().GetValue("WEB_SITE_DOMAIN")
	WebConfig := NewConfigService().GetValue("SSL_OPEN")
	if webSiteDomain != "" {
		if WebConfig == "1" {
			http = "https://"
		} else {
			http = "http://"
		}
	}

	if getId, ok := id.(string); ok {
		// json字符串
		if strings.Contains(getId, "{") {
			var jsonData []map[string]interface{}
			err := json.Unmarshal([]byte(getId), &jsonData)
			if err == nil {
				for _, v := range jsonData {
					path = v["url"].(string)
					if strings.Contains(path, "//") {
						paths = append(paths, v["url"].(string))
					} else {
						if strings.Contains(path, "./") {
							path = strings.Replace(path, "./web/app/", "/", -1)
						}
						if path != "" {
							path = http + webSiteDomain + path
						}
						paths = append(paths, path)
					}
				}
			}
		}
	}

	return paths
}
