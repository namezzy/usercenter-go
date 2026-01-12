package controller

import (
	"net/http"
	"user-center/internal/common"
	"user-center/internal/model"
	"user-center/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// Register 用户注册
func (ctrl *UserController) Register(c *gin.Context) {
	var req model.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, common.NewErrorResponse(common.PARAMS_ERROR))
		return
	}

	userId, err := ctrl.userService.UserRegister(req.UserAccount, req.UserPassword, req.CheckPassword, req.PlanetCode)
	if err != nil {
		c.JSON(http.StatusOK, common.NewErrorResponseWithDesc(common.PARAMS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(userId))
}

// Login 用户登录
func (ctrl *UserController) Login(c *gin.Context) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, common.NewErrorResponse(common.PARAMS_ERROR))
		return
	}

	user, err := ctrl.userService.UserLogin(req.UserAccount, req.UserPassword)
	if err != nil {
		c.JSON(http.StatusOK, common.NewErrorResponseWithDesc(common.PARAMS_ERROR, err.Error()))
		return
	}

	// 记录登录态
	session := sessions.Default(c)
	session.Set(common.UserLoginState, user)
	session.Save()

	c.JSON(http.StatusOK, common.NewSuccessResponse(user))
}

// Logout 用户注销
func (ctrl *UserController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(common.UserLoginState)
	session.Save()

	c.JSON(http.StatusOK, common.NewSuccessResponse(1))
}

// GetCurrentUser 获取当前用户
func (ctrl *UserController) GetCurrentUser(c *gin.Context) {
	session := sessions.Default(c)
	userObj := session.Get(common.UserLoginState)

	if userObj == nil {
		c.JSON(http.StatusOK, common.NewErrorResponse(common.NOT_LOGIN))
		return
	}

	currentUser, ok := userObj.(*model.SafetyUser)
	if !ok {
		c.JSON(http.StatusOK, common.NewErrorResponse(common.SYSTEM_ERROR))
		return
	}

	// 重新从数据库获取最新用户信息
	user, err := ctrl.userService.GetUserByID(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusOK, common.NewErrorResponse(common.SYSTEM_ERROR))
		return
	}

	safetyUser := ctrl.userService.GetSafetyUser(user)
	c.JSON(http.StatusOK, common.NewSuccessResponse(safetyUser))
}

// SearchUsers 搜索用户
func (ctrl *UserController) SearchUsers(c *gin.Context) {
	username := c.Query("username")

	users, err := ctrl.userService.SearchUsers(username)
	if err != nil {
		c.JSON(http.StatusOK, common.NewErrorResponse(common.SYSTEM_ERROR))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(users))
}

// DeleteUser 删除用户
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, common.NewErrorResponse(common.PARAMS_ERROR))
		return
	}

	err := ctrl.userService.DeleteUser(req.ID)
	if err != nil {
		c.JSON(http.StatusOK, common.NewErrorResponseWithDesc(common.PARAMS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(true))
}
