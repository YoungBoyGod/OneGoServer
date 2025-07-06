package system

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/system/v1"
)

func (c *ControllerV1) GetSystemStatus(ctx context.Context, req *v1.GetSystemStatusReq) (res *v1.GetSystemStatusRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
