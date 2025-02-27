package upload

import (
	"reflect"
	"time"

	"github.com/quarkcloudio/quark-go/v3"
	"github.com/quarkcloudio/quark-go/v3/model"
	"github.com/quarkcloudio/quark-go/v3/service"
	"github.com/quarkcloudio/quark-go/v3/template/tool/upload"
)

type File struct {
	upload.Template
}

// 初始化
func (p *File) Init(ctx *quark.Context) interface{} {

	// 限制文件大小
	p.LimitSize = 1024 * 1024 * 1024 * 2

	// 限制文件类型
	p.LimitType = []string{
		"image/png",
		"image/gif",
		"image/jpeg",
		"video/mp4",
		"video/mpeg",
		"application/x-xls",
		"application/x-ppt",
		"application/msword",
		"application/zip",
		"application/pdf",
		"application/vnd.ms-excel",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	}

	// 设置文件上传路径
	p.SavePath = "./web/app/storage/files/" + time.Now().Format("20060102") + "/"

	return p
}

// 上传前回调
func (p *File) BeforeHandle(ctx *quark.Context, fileSystem *quark.FileSystem) (*quark.FileSystem, *quark.FileInfo, error) {
	fileHash, err := fileSystem.GetFileHash()
	if err != nil {
		return fileSystem, nil, err
	}

	getFileInfo, _ := service.NewFileService().GetInfoByHash(fileHash)
	if getFileInfo.Id != 0 {
		fileInfo := &quark.FileInfo{
			Name: getFileInfo.Name,
			Size: getFileInfo.Size,
			Ext:  getFileInfo.Ext,
			Path: getFileInfo.Path,
			Url:  getFileInfo.Url,
			Hash: getFileInfo.Hash,
		}

		return fileSystem, fileInfo, err
	}

	return fileSystem, nil, err
}

// 上传完成后回调
func (p *File) AfterHandle(ctx *quark.Context, result *quark.FileInfo) error {
	driver := reflect.
		ValueOf(ctx.Template).
		Elem().
		FieldByName("Driver").String()

	// 重写url
	if driver == quark.LocalStorage {
		result.Url = service.NewFileService().GetPath(result.Url)
	}

	// 插入数据库
	id, err := service.NewFileService().InsertGetId(model.File{
		Name:   result.Name,
		Size:   result.Size,
		Ext:    result.Ext,
		Path:   result.Path,
		Url:    result.Url,
		Hash:   result.Hash,
		Status: 1,
	})
	if err != nil {
		return ctx.JSONError(err.Error())
	}

	return ctx.JSONOk("上传成功", map[string]interface{}{
		"id":          id,
		"contentType": result.ContentType,
		"ext":         result.Ext,
		"hash":        result.Hash,
		"name":        result.Name,
		"path":        result.Path,
		"size":        result.Size,
		"url":         result.Url,
	})
}
