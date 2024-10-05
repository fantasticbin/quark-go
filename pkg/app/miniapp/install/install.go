package install

import (
	"github.com/quarkcloudio/quark-go/v3/pkg/app/miniapp/model"
	"github.com/quarkcloudio/quark-go/v3/pkg/dal/db"
	"gorm.io/gorm"
)

// 执行安装操作
func Handle() {

	// 迁移数据
	db.Client.AutoMigrate(
		&model.User{},
	)

	// 如果用户不存在，初始化数据库数据
	userInfo, err := (&model.User{}).GetInfoById(1)
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	if userInfo.Id == 0 {
		// 数据填充
		(&model.User{}).Seeder()
	}
}
