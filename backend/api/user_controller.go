package api

import (
	"net/http"
	"react-mng2-backend/models"
	"react-mng2-backend/services"
	"react-mng2-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: &services.UserService{},
	}
}

// GetList 获取用户列表
func (ctrl *UserController) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	users, total, err := ctrl.userService.GetUserList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, utils.Error("获取用户列表失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(utils.NewPageData(users, total, page, pageSize)))
}

// GetDetail 获取用户详情
func (ctrl *UserController) GetDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	user, err := ctrl.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, utils.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Success(user))
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Realname string `json:"realname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// Create 创建用户
func (ctrl *UserController) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("参数错误"))
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Realname: req.Realname,
		Email:    req.Email,
		Status:   1,
	}

	if err := ctrl.userService.CreateUser(user); err != nil {
		c.JSON(http.StatusOK, utils.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Success(user))
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username"`
	Realname string `json:"realname"`
	Email    string `json:"email"`
	Status   *int   `json:"status"`
}

// Update 更新用户
func (ctrl *UserController) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("参数错误"))
		return
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Realname != "" {
		updates["realname"] = req.Realname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := ctrl.userService.UpdateUser(uint(id), updates); err != nil {
		c.JSON(http.StatusOK, utils.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Success(nil))
}

// Delete 删除用户
func (ctrl *UserController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := ctrl.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusOK, utils.Error("删除用户失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(nil))
}
