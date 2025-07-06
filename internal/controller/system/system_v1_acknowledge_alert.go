package system

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/system/v1"
)

func (c *ControllerV1) AcknowledgeAlert(ctx context.Context, req *v1.AcknowledgeAlertReq) (res *v1.AcknowledgeAlertRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
