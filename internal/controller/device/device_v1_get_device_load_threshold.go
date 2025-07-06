package device

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/device/v1"
)

func (c *ControllerV1) GetDeviceLoadThreshold(ctx context.Context, req *v1.GetDeviceLoadThresholdReq) (res *v1.GetDeviceLoadThresholdRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
