package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/user/v1"
)

func (c *ControllerV1) CheckUserPermission(ctx context.Context, req *v1.CheckUserPermissionReq) (res *v1.CheckUserPermissionRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
