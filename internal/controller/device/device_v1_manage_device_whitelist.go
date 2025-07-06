package device

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/device/v1"
)

func (c *ControllerV1) ManageDeviceWhitelist(ctx context.Context, req *v1.ManageDeviceWhitelistReq) (res *v1.ManageDeviceWhitelistRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
