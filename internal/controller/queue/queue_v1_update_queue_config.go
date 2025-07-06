package queue

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/queue/v1"
)

func (c *ControllerV1) UpdateQueueConfig(ctx context.Context, req *v1.UpdateQueueConfigReq) (res *v1.UpdateQueueConfigRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
