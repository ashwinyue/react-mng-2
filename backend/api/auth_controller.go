package api

import (
	"net/http"
	"react-mng2-backend/services"
	"react-mng2-backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	userService *services.UserService
}

// NewAuthController 创建认证控制器
func NewAuthController() *AuthController {
	return &AuthController{
		userService: &services.UserService{},
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("参数错误"))
		return
	}

	// 查询用户
	user, err := ctrl.userService.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusOK, utils.Error("用户名或密码错误"))
		return
	}

	// 验证密码
	if !ctrl.userService.VerifyPassword(user, req.Password) {
		c.JSON(http.StatusOK, utils.Error("用户名或密码错误"))
		return
	}

	// 生成 token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusOK, utils.Error("生成 token 失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"realname": user.Realname,
			"email":    user.Email,
		},
	}))
}

// Logout 用户登出
func (ctrl *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, utils.Success(nil))
}

// GetProfile 获取当前用户信息
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := ctrl.userService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, utils.Error("获取用户信息失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
		"realname": user.Realname,
		"email":    user.Email,
		"status":   user.Status,
	}))
}
