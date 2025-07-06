package device

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/device/v1"
)

func (c *ControllerV1) UpdateDeviceInfo(ctx context.Context, req *v1.UpdateDeviceInfoReq) (res *v1.UpdateDeviceInfoRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
