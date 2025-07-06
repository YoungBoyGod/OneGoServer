package device

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/device/v1"
)

func (c *ControllerV1) GetDeviceLoadStatus(ctx context.Context, req *v1.GetDeviceLoadStatusReq) (res *v1.GetDeviceLoadStatusRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
