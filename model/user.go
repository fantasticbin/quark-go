package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/quarkcloudio/quark-go/v3/dal/db"
	"github.com/quarkcloudio/quark-go/v3/utils/datetime"
	"github.com/quarkcloudio/quark-go/v3/utils/hash"
	"gorm.io/gorm"
)

// 字段
type User struct {
	Id            int               `json:"id" gorm:"autoIncrement"`
	Username      string            `json:"username" gorm:"size:20;index:admins_username_unique,unique;not null"`
	Nickname      string            `json:"nickname" gorm:"size:200;not null"`
	Sex           int               `json:"sex" gorm:"size:4;not null;default:1"`
	Email         string            `json:"email" gorm:"size:50;index:admins_email_unique,unique;not null"`
	Phone         string            `json:"phone" gorm:"size:11;index:admins_phone_unique,unique;not null"`
	Password      string            `json:"password" gorm:"size:255;not null"`
	Avatar        string            `json:"avatar" gorm:"size:1000"`
	DepartmentId  int               `json:"department_id" gorm:"size:11;not null;default:0"`
	PositionIds   string            `json:"position_ids" gorm:"size:1000;default:null"`
	LastLoginIp   string            `json:"last_login_ip" gorm:"size:255"`
	LastLoginTime datetime.Datetime `json:"last_login_time"`
	WxOpenid      string            `json:"wx_openid" gorm:"size:255"`
	WxUnionid     string            `json:"wx_unionid" gorm:"size:255"`
	Status        int               `json:"status" gorm:"size:1;not null;default:1"`
	CreatedAt     datetime.Datetime `json:"created_at"`
	UpdatedAt     datetime.Datetime `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `json:"deleted_at"`
}

// 管理员JWT结构体
type UserClaims struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Sex       int    `json:"sex"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	GuardName string `json:"guard_name"`
	jwt.RegisteredClaims
}

// 管理员Seeder
func (model *User) Seeder() {
	seeders := []User{
		{Username: "administrator", Nickname: "超级管理员", Email: "admin@yourweb.com", Phone: "10086", Password: hash.Make("123456"), Sex: 1, DepartmentId: 1, Status: 1, LastLoginTime: datetime.Now()},
	}

	db.Client.Create(&seeders)
}

// 获取管理员JWT信息
func (model *User) GetAdminClaims(adminInfo User) (adminClaims *UserClaims) {
	adminClaims = &UserClaims{
		adminInfo.Id,
		adminInfo.Username,
		adminInfo.Nickname,
		adminInfo.Sex,
		adminInfo.Email,
		adminInfo.Phone,
		adminInfo.Avatar,
		"admin",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 过期时间，默认24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 颁发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                     // 不早于时间
			Issuer:    "QuarkGo",                                          // 颁发人
			Subject:   "Admin Token",                                      // 主题信息
		},
	}

	return adminClaims
}

// 获取普通用户JWT信息
func (model *User) GetUserClaims(adminInfo *User) (adminClaims *UserClaims) {
	adminClaims = &UserClaims{
		adminInfo.Id,
		adminInfo.Username,
		adminInfo.Nickname,
		adminInfo.Sex,
		adminInfo.Email,
		adminInfo.Phone,
		adminInfo.Avatar,
		"user",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 过期时间，默认24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 颁发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                     // 不早于时间
			Issuer:    "QuarkGo",                                          // 颁发人
			Subject:   "User Token",                                       // 主题信息
		},
	}

	return adminClaims
}

// 获取当前认证的用户信息，默认参数为tokenString
func (model *User) GetAuthUser(appKey string, tokenString string) (adminClaims *UserClaims, Error error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(appKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("token格式错误")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token已过期")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token未生效")
			} else {
				return nil, err
			}
		}
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token不可用")
}

// 通过ID获取管理员信息
func (model *User) GetInfoById(id interface{}) (admin User, Error error) {
	err := db.Client.Where("status = ?", 1).Where("id = ?", id).First(&admin).Error

	return admin, err
}

// 通过用户名获取管理员信息
func (model *User) GetInfoByUsername(username string) (admin User, Error error) {
	err := db.Client.Where("status = ?", 1).Where("username = ?", username).First(&admin).Error
	if admin.Avatar != "" {
		admin.Avatar = (&Picture{}).GetPath(admin.Avatar) // 获取头像地址
	}

	return admin, err
}

// 通过ID获取管理员拥有的菜单列表
func (model *User) GetMenuListById(id interface{}) (menuList interface{}, Error error) {

	return (&Menu{}).GetListByAdminId(id.(int))
}

// 更新最后一次登录数据
func (model *User) UpdateLastLogin(uid int, lastLoginIp string, lastLoginTime datetime.Datetime) error {
	data := User{
		LastLoginIp:   lastLoginIp,
		LastLoginTime: lastLoginTime,
	}

	return db.Client.
		Where("id = ?", uid).
		Updates(&data).Error
}
