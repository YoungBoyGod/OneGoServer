package v1

import (
	"time"
)

// CommonRes 通用响应结构
type CommonRes struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PageRes 分页响应结构
type PageRes struct {
	CommonRes
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// UserRegisterReq 用户注册请求
type UserRegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
}

// UserLoginReq 用户登录请求
type UserLoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateReq 用户更新请求
type UserUpdateReq struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Avatar   string `json:"avatar" binding:"omitempty,url"`
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UserResp 用户响应
type UserResp struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Avatar      string     `json:"avatar"`
	Status      int        `json:"status"`
	Role        int        `json:"role"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserLoginResp 用户登录响应
type UserLoginResp struct {
	User  *UserResp `json:"user"`
	Token string    `json:"token"`
}

// UserListReq 用户列表请求
type UserListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Username string `form:"username"`
	Email    string `form:"email"`
	Status   *int   `form:"status"`
	Role     *int   `form:"role"`
}
