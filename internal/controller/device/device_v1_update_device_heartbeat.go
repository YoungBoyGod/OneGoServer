package device

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/device/v1"
)

func (c *ControllerV1) UpdateDeviceHeartbeat(ctx context.Context, req *v1.UpdateDeviceHeartbeatReq) (res *v1.UpdateDeviceHeartbeatRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
