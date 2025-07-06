package system

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/system/v1"
)

func (c *ControllerV1) ResolveAlert(ctx context.Context, req *v1.ResolveAlertReq) (res *v1.ResolveAlertRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
