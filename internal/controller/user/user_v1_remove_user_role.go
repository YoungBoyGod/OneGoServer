package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/user/v1"
)

func (c *ControllerV1) RemoveUserRole(ctx context.Context, req *v1.RemoveUserRoleReq) (res *v1.RemoveUserRoleRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
