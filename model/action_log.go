package model

import (
	"github.com/quarkcloudio/quark-go/v3/utils/datetime"
)

// 字段
type ActionLog struct {
	Id        int               `json:"id" gorm:"autoIncrement"`
	ObjectId  int               `json:"object_id" gorm:"size:11;not null"`
	Username  string            `json:"username" gorm:"<-:false"`
	Url       string            `json:"url" gorm:"size:500;not null"`
	Remark    string            `json:"remark" gorm:"size:255;not null"`
	Ip        string            `json:"ip" gorm:"size:100;not null"`
	Type      string            `json:"type" gorm:"size:100;not null"`
	Status    int               `json:"status" gorm:"size:1;not null;default:1"`
	CreatedAt datetime.Datetime `json:"created_at"`
	UpdatedAt datetime.Datetime `json:"updated_at"`
}
