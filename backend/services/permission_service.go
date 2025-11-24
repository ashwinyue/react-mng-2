package services

import (
	"react-mng2-backend/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"gorm.io/gorm"
)

// PermissionService 权限服务
type PermissionService struct {
}

// PermissionTree 权限树结构
type PermissionTree struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Code        string            `json:"code"`
	ParentCode  string            `json:"parent_code"`
	Path        string            `json:"path"`
	Type        int               `json:"type"`
	Sort        int               `json:"sort"`
	Description string            `json:"description"`
	Children    []*PermissionTree `json:"children,omitempty"`
	HasChildren bool              `json:"has_children"`
}

// RolePermissionData 角色权限数据
type RolePermissionData struct {
	RoleID       uint             `json:"role_id"`
	PermissionTrees []*PermissionTree `json:"permission_trees"`
}

// StringSlice 自定义类型用于JSON序列化
type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// GetAllPermissions 获取所有权限
func (s *PermissionService) GetAllPermissions() ([]*models.Permission, error) {
	var permissions []*models.Permission
	if err := models.DB.Preload("Roles").Order("type ASC, sort ASC, id ASC").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetPermissionTrees 获取权限树结构
func (s *PermissionService) GetPermissionTrees() ([]*PermissionTree, error) {
	permissions, err := s.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	// 构建权限树
	treeMap := make(map[string]*PermissionTree)
	var roots []*PermissionTree

	// 先创建所有节点
	for _, perm := range permissions {
		treeMap[perm.Code] = &PermissionTree{
			ID:          perm.ID,
			Name:        perm.Name,
			Code:        perm.Code,
			ParentCode:  perm.ParentCode,
			Path:        perm.Path,
			Type:        perm.Type,
			Sort:        perm.Sort,
			Description: perm.Description,
			Children:    make([]*PermissionTree, 0),
			HasChildren: false,
		}
	}

	// 构建父子关系
	for _, perm := range permissions {
		node := treeMap[perm.Code]
		if perm.ParentCode == "" {
			// 根节点
			roots = append(roots, node)
		} else {
			// 子节点
			if parent, exists := treeMap[perm.ParentCode]; exists {
				parent.Children = append(parent.Children, node)
				parent.HasChildren = true
			}
		}
	}

	// 对根节点排序
	sort.Slice(roots, func(i, j int) bool {
		if roots[i].Type == roots[j].Type {
			return roots[i].Sort < roots[j].Sort
		}
		return roots[i].Type < roots[j].Type
	})

	// 对子节点排序
	var sortChildren func([]*PermissionTree)
	sortChildren = func(nodes []*PermissionTree) {
		for i := range nodes {
			if len(nodes[i].Children) > 0 {
				sort.Slice(nodes[i].Children, func(j, k int) bool {
					if nodes[i].Children[j].Type == nodes[i].Children[k].Type {
						return nodes[i].Children[j].Sort < nodes[i].Children[k].Sort
					}
					return nodes[i].Children[j].Type < nodes[i].Children[k].Type
				})
				sortChildren(nodes[i].Children)
			}
		}
	}
	sortChildren(roots)

	return roots, nil
}

// GetRolePermissions 获取角色的权限
func (s *PermissionService) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	var role models.Role
	if err := models.DB.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

// AssignPermissionsToRole 为角色分配权限
func (s *PermissionService) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	var role models.Role
	if err := models.DB.First(&role, roleID).Error; err != nil {
		return err
	}

	// 清空现有权限
	if err := models.DB.Model(&role).Association("Permissions").Clear(); err != nil {
		return err
	}

	// 分配新权限
	var permissions []*models.Permission
	if len(permissionIDs) > 0 {
		if err := models.DB.Find(&permissions, permissionIDs).Error; err != nil {
			return err
		}
		if err := models.DB.Model(&role).Association("Permissions").Append(permissions); err != nil {
			return err
		}
	}

	return nil
}

// CheckPermission 检查用户是否具有特定权限
func (s *PermissionService) CheckPermission(userID uint, permissionCode string) bool {
	var user models.User
	if err := models.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return false
	}

	if user.Role == nil {
		return false
	}

	for _, perm := range user.Role.Permissions {
		if perm.Code == permissionCode {
			return true
		}
	}

	return false
}

// GetUserPermissions 获取用户的所有权限代码
func (s *PermissionService) GetUserPermissions(userID uint) []string {
	var user models.User
	if err := models.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return nil
	}

	if user.Role == nil {
		return nil
	}

	var codes []string
	for _, perm := range user.Role.Permissions {
		codes = append(codes, perm.Code)
	}

	return codes
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(permission *models.Permission) error {
	// 检查权限代码是否已存在
	var existingPerm models.Permission
	if err := models.DB.Where("code = ?", permission.Code).First(&existingPerm).Error; err == nil {
		return fmt.Errorf("permission with code '%s' already exists", permission.Code)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return models.DB.Create(permission).Error
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(permission *models.Permission) error {
	// 检查记录是否存在
	var existingPerm models.Permission
	if err := models.DB.First(&existingPerm, permission.ID).Error; err != nil {
		return err
	}

	// 检查代码冲突（排除自身）
	if permission.Code != existingPerm.Code {
		var conflictPerm models.Permission
		if err := models.DB.Where("code = ? AND id != ?", permission.Code, permission.ID).First(&conflictPerm).Error; err == nil {
			return fmt.Errorf("permission with code '%s' already exists", permission.Code)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	return models.DB.Save(permission).Error
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(id uint) error {
	// 开启事务
	tx := models.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查权限是否存在
	var permission models.Permission
	if err := tx.First(&permission, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除角色权限关联
	if err := tx.Model(&permission).Association("Roles").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// 执行删除
	if err := tx.Delete(&permission).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetPermissionByID 根据ID获取权限
func (s *PermissionService) GetPermissionByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	if err := models.DB.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetPermissionByCode 根据代码获取权限
func (s *PermissionService) GetPermissionByCode(code string) (*models.Permission, error) {
	var permission models.Permission
	if err := models.DB.Where("code = ?", code).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetPermissionTreesWithRole 获取角色权限树结构
func (s *PermissionService) GetPermissionTreesWithRole(roleID uint) (*RolePermissionData, error) {
	// 获取权限树
	trees, err := s.GetPermissionTrees()
	if err != nil {
		return nil, err
	}

	// 获取角色权限
	rolePerms, err := s.GetRolePermissions(roleID)
	if err != nil {
		return nil, err
	}

	// 创建权限代码集合
	rolePermCodes := make(map[uint]bool)
	for _, perm := range rolePerms {
		rolePermCodes[perm.ID] = true
	}

	// 标记选中的权限
	var markSelected func([]*PermissionTree)
	markSelected = func(nodes []*PermissionTree) {
		for _, node := range nodes {
			if rolePermCodes[node.ID] {
				// 标记为选中状态（这里可以添加具体的状态标记逻辑）
				// 例如：node.Selected = true
			}
			if len(node.Children) > 0 {
				markSelected(node.Children)
			}
		}
	}
	markSelected(trees)

	return &RolePermissionData{
		RoleID:         roleID,
		PermissionTrees: trees,
	}, nil
}