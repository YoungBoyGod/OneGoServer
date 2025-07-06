package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/user/v1"
)

func (c *ControllerV1) RevokeUserSession(ctx context.Context, req *v1.RevokeUserSessionReq) (res *v1.RevokeUserSessionRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
