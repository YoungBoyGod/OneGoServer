package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/user/v1"
)

func (c *ControllerV1) GetUserSecurityLog(ctx context.Context, req *v1.GetUserSecurityLogReq) (res *v1.GetUserSecurityLogRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
