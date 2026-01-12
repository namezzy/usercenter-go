package model

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	UserAccount   string `json:"userAccount" binding:"required"`
	UserPassword  string `json:"userPassword" binding:"required"`
	CheckPassword string `json:"checkPassword" binding:"required"`
	PlanetCode    string `json:"planetCode" binding:"required"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	UserAccount  string `json:"userAccount" binding:"required"`
	UserPassword string `json:"userPassword" binding:"required"`
}
