package model

import (
	"github.com/quarkcloudio/quark-go/v3/utils/datetime"
)

// 字段
type Picture struct {
	Id                int               `json:"id" gorm:"autoIncrement"`
	ObjType           string            `json:"obj_type" gorm:"size:255"`
	ObjId             int               `json:"obj_id" gorm:"size:11;default:0"`
	PictureCategoryId int               `json:"picture_category_id" gorm:"size:11;default:0"`
	Sort              int               `json:"sort" gorm:"size:11;default:0"`
	Name              string            `json:"name" gorm:"size:255;not null"`
	Size              int64             `json:"size" gorm:"size:20;default:0"`
	Width             int               `json:"width" gorm:"size:11;default:0"`
	Height            int               `json:"height" gorm:"size:11;default:0"`
	Ext               string            `json:"ext" gorm:"size:255"`
	Path              string            `json:"path" gorm:"size:255;not null"`
	Url               string            `json:"url" gorm:"size:255;not null"`
	Hash              string            `json:"hash" gorm:"size:255;not null"`
	Status            int               `json:"status" gorm:"size:1;not null;default:1"`
	CreatedAt         datetime.Datetime `json:"created_at"`
	UpdatedAt         datetime.Datetime `json:"updated_at"`
}
