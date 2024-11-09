package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/model"
	"github.com/xuri/excelize/v2"
)

type FileService struct{}

// 初始化
func NewFileService() *FileService {
	return &FileService{}
}

// 插入数据并返回ID
func (p *FileService) InsertGetId(data model.File) (id int, Error error) {
	err := db.Client.Create(&data).Error

	return data.Id, err
}

// 根据hash查询文件信息
func (p *FileService) GetInfoByHash(hash string) (file model.File, Error error) {
	err := db.Client.Where("status = ?", 1).Where("hash = ?", hash).First(&file).Error

	return file, err
}

// 获取文件路径
func (p *FileService) GetPath(id interface{}) string {
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
			// 如果为map
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
			return http + webSiteDomain + path
		}
	}

	file := model.File{}
	db.Client.Where("id", id).Where("status", 1).First(&file)
	if file.Id != 0 {
		path = file.Url
		if strings.Contains(path, "//") {
			return path
		}
		if strings.Contains(path, "./") {
			path = strings.Replace(path, "./web/app/", "/", -1)
		}
		if path != "" {
			return http + webSiteDomain + path
		}
	}

	return ""
}

// 获取多文件路径
func (p *FileService) GetPaths(id interface{}) []string {
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

// 获取Excel文件数据
func (p *FileService) GetExcelData(fileId int) (data [][]interface{}, Error error) {
	file := model.File{}
	err := db.Client.Where("id", fileId).Where("status", 1).First(&file).Error
	if err != nil {
		return data, err
	}
	if file.Id == 0 {
		return data, errors.New("参数错误！")
	}

	f, err := excelize.OpenFile(file.Path)
	if err != nil {
		return data, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return data, err
	}

	for _, row := range rows {
		getRows := []interface{}{}
		for _, colCell := range row {
			getRows = append(getRows, colCell)
		}
		data = append(data, getRows)
	}

	return data, err
}
