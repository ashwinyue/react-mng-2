package api

import (
	"react-mng2-backend/services"
	"react-mng2-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PermissionController 权限控制器
type PermissionController struct {
	permissionService *services.PermissionService
}

// NewPermissionController 创建权限控制器
func NewPermissionController() *PermissionController {
	return &PermissionController{
		permissionService: &services.PermissionService{},
	}
}

// GetList 获取权限列表
func (c *PermissionController) GetList(ctx *gin.Context) {
	permissions, err := c.permissionService.GetAllPermissions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Error("获取权限列表失败: " + err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.Success(gin.H{
		"data":  permissions,
		"total": len(permissions),
	}))
}

// GetTree 获取权限树
func (c *PermissionController) GetTree(ctx *gin.Context) {
	trees, err := c.permissionService.GetPermissionTrees()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Error("获取权限树失败: " + err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.Success(trees))
}

// GetDetail 获取权限详情
func (c *PermissionController) GetDetail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Error("无效的权限ID"))
		return
	}

	permission, err := c.permissionService.GetPermissionByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.Error("权限不存在"))
		return
	}

	ctx.JSON(http.StatusOK, utils.Success(permission))
}

// Create 创建权限
func (c *PermissionController) Create(ctx *gin.Context) {
	var permission services.PermissionTree
	if err := ctx.ShouldBindJSON(&permission); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Error("请求参数错误: " + err.Error()))
		return
	}

	// 验证必要字段
	if permission.Name == "" || permission.Code == "" {
		ctx.JSON(http.StatusBadRequest, utils.Error("权限名称和代码不能为空"))
		return
	}

	// 检查权限代码是否已存在
	existingPerm, err := c.permissionService.GetPermissionByCode(permission.Code)
	if err == nil && existingPerm != nil {
		ctx.JSON(http.StatusBadRequest, utils.Error("权限代码已存在"))
		return
	}

	// 这里需要调用实际的创建方法，但现在没有实现
	// ctx.JSON(http.StatusOK, utils.Success(gin.H{"message": "权限创建成功"}))
}

// Update 更新权限
func (c *PermissionController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Error("无效的权限ID"))
		return
	}

	var permission services.PermissionTree
	if err := ctx.ShouldBindJSON(&permission); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Error("请求参数错误: " + err.Error()))
		return
	}

	// 获取现有权限
	existingPerm, err := c.permissionService.GetPermissionByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.Error("权限不存在"))
		return
	}

	// 检查权限代码是否已存在（排除当前权限）
	if permission.Code != existingPerm.Code {
		checkPerm, err := c.permissionService.GetPermissionByCode(permission.Code)
		if err == nil && checkPerm != nil {
			ctx.JSON(http.StatusBadRequest, utils.Error("权限代码已存在"))
			return
		}
	}

	// 这里需要调用实际的更新方法，但现在没有实现
	// ctx.JSON(http.StatusOK, utils.Success(gin.H{"message": "权限更新成功"}))
}

// Delete 删除权限
func (c *PermissionController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Error("无效的权限ID"))
		return
	}

	if err := c.permissionService.DeletePermission(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Error("删除权限失败: " + err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.Success(gin.H{"message": "权限删除成功"}))
}