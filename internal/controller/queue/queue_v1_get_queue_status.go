package queue

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/queue/v1"
)

func (c *ControllerV1) GetQueueStatus(ctx context.Context, req *v1.GetQueueStatusReq) (res *v1.GetQueueStatusRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
