package dto

import "github.com/golang-jwt/jwt/v4"

// 管理员JWT结构体
type UserClaims struct {
	jwt.RegisteredClaims
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Sex       int    `json:"sex"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	GuardName string `json:"guard_name"`
}
