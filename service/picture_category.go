package service

import (
	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/model"
)

type PictureCategoryService struct{}

// 初始化
func NewPictureCategoryService() *PictureCategoryService {
	return &PictureCategoryService{}
}

// 获取列表
func (p *PictureCategoryService) GetAuthList(appKey string, tokenString string) (list []model.PictureCategory, Error error) {
	categorys := []model.PictureCategory{}

	adminInfo, err := NewUserService().GetAuthUser(appKey, tokenString)
	if err != nil {
		return categorys, err
	}

	err = db.Client.
		Where("obj_type = ?", "ADMIN").
		Where("obj_id", adminInfo.Id).
		Find(&categorys).Error
	if err != nil {
		return categorys, err
	}

	return categorys, nil
}
