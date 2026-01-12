package middleware

import (
	"net/http"
	"user-center/internal/common"
	"user-center/internal/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userObj := session.Get(common.UserLoginState)
		
		if userObj == nil {
			c.JSON(http.StatusOK, common.NewErrorResponse(common.NOT_LOGIN))
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// AdminAuthMiddleware 管理员权限中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userObj := session.Get(common.UserLoginState)
		
		if userObj == nil {
			c.JSON(http.StatusOK, common.NewErrorResponse(common.NOT_LOGIN))
			c.Abort()
			return
		}
		
		user, ok := userObj.(*model.SafetyUser)
		if !ok || user.UserRole != common.AdminRole {
			c.JSON(http.StatusOK, common.NewErrorResponse(common.NO_AUTH))
			c.Abort()
			return
		}
		
		c.Next()
	}
}
