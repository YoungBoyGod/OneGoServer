package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/user/v1"
)

func (c *ControllerV1) GetUserRoles(ctx context.Context, req *v1.GetUserRolesReq) (res *v1.GetUserRolesRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
