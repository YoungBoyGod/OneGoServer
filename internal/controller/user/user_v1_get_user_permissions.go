package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/user/v1"
)

func (c *ControllerV1) GetUserPermissions(ctx context.Context, req *v1.GetUserPermissionsReq) (res *v1.GetUserPermissionsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
