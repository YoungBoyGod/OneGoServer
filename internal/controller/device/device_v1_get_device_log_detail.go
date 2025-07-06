package device

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/device/v1"
)

func (c *ControllerV1) GetDeviceLogDetail(ctx context.Context, req *v1.GetDeviceLogDetailReq) (res *v1.GetDeviceLogDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
