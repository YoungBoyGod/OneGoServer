// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"OneGfServer/api/user/v1"
)

type IUserV1 interface {
	RegisterUser(ctx context.Context, req *v1.RegisterUserReq) (res *v1.RegisterUserRes, err error)
	LoginUser(ctx context.Context, req *v1.LoginUserReq) (res *v1.LoginUserRes, err error)
	LogoutUser(ctx context.Context, req *v1.LogoutUserReq) (res *v1.LogoutUserRes, err error)
	GetUserList(ctx context.Context, req *v1.GetUserListReq) (res *v1.GetUserListRes, err error)
	GetUserDetail(ctx context.Context, req *v1.GetUserDetailReq) (res *v1.GetUserDetailRes, err error)
	UpdateUser(ctx context.Context, req *v1.UpdateUserReq) (res *v1.UpdateUserRes, err error)
	GetUserRoles(ctx context.Context, req *v1.GetUserRolesReq) (res *v1.GetUserRolesRes, err error)
	AssignUserRole(ctx context.Context, req *v1.AssignUserRoleReq) (res *v1.AssignUserRoleRes, err error)
	RemoveUserRole(ctx context.Context, req *v1.RemoveUserRoleReq) (res *v1.RemoveUserRoleRes, err error)
	GetUserPermissions(ctx context.Context, req *v1.GetUserPermissionsReq) (res *v1.GetUserPermissionsRes, err error)
	CheckUserPermission(ctx context.Context, req *v1.CheckUserPermissionReq) (res *v1.CheckUserPermissionRes, err error)
	GetUserSessions(ctx context.Context, req *v1.GetUserSessionsReq) (res *v1.GetUserSessionsRes, err error)
	RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error)
	RevokeUserSession(ctx context.Context, req *v1.RevokeUserSessionReq) (res *v1.RevokeUserSessionRes, err error)
	GetUserActivity(ctx context.Context, req *v1.GetUserActivityReq) (res *v1.GetUserActivityRes, err error)
	ChangePassword(ctx context.Context, req *v1.ChangePasswordReq) (res *v1.ChangePasswordRes, err error)
	ResetPassword(ctx context.Context, req *v1.ResetPasswordReq) (res *v1.ResetPasswordRes, err error)
	GetUserSecurityLog(ctx context.Context, req *v1.GetUserSecurityLogReq) (res *v1.GetUserSecurityLogRes, err error)
	UpdateUserSecuritySettings(ctx context.Context, req *v1.UpdateUserSecuritySettingsReq) (res *v1.UpdateUserSecuritySettingsRes, err error)
	GetUserStatistics(ctx context.Context, req *v1.GetUserStatisticsReq) (res *v1.GetUserStatisticsRes, err error)
	GetUserBehaviorAnalysis(ctx context.Context, req *v1.GetUserBehaviorAnalysisReq) (res *v1.GetUserBehaviorAnalysisRes, err error)
}
