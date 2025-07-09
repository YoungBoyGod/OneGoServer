// package user

// import (
// 	"context"
// 	"fmt"
// 	"strings"

// 	"OneGfServer/internal/model/user"
// )

// // ===============================
// // 用户权限管理相关业务逻辑
// // ===============================

// // CalculateUserPermissions 计算用户权限
// func (s *sUser) CalculateUserPermissions(ctx context.Context, input *user.CalculateUserPermissionsInput) (*user.CalculateUserPermissionsOutput, error) {
// 	permissionSet := make(map[string]bool)

// 	for _, role := range input.UserRoles {
// 		if permissions, ok := role["permissions"].([]string); ok {
// 			for _, permission := range permissions {
// 				permissionSet[permission] = true
// 			}
// 		}
// 	}

// 	// 转换为数组
// 	var permissions []string
// 	for permission := range permissionSet {
// 		permissions = append(permissions, permission)
// 	}

// 	return &user.CalculateUserPermissionsOutput{
// 		Permissions: permissions,
// 	}, nil
// }

// // CheckUserPermission 检查用户权限
// func (s *sUser) CheckUserPermission(ctx context.Context, input *user.CheckUserPermissionInput) (*user.CheckUserPermissionOutput, error) {
// 	// 检查直接权限
// 	for _, permission := range input.UserPermissions {
// 		if permission == input.RequiredPermission {
// 			return &user.CheckUserPermissionOutput{
// 				HasPermission: true,
// 			}, nil
// 		}

// 		// 检查通配符权限
// 		if strings.HasSuffix(permission, "*") {
// 			prefix := strings.TrimSuffix(permission, "*")
// 			if strings.HasPrefix(input.RequiredPermission, prefix) {
// 				return &user.CheckUserPermissionOutput{
// 					HasPermission: true,
// 				}, nil
// 			}
// 		}
// 	}

// 	// 检查资源特定权限
// 	resourcePermission := fmt.Sprintf("%s:%s", input.Resource, input.RequiredPermission)
// 	for _, permission := range input.UserPermissions {
// 		if permission == resourcePermission {
// 			return &user.CheckUserPermissionOutput{
// 				HasPermission: true,
// 			}, nil
// 		}
// 	}

// 	return &user.CheckUserPermissionOutput{
// 		HasPermission: false,
// 	}, nil
// }
