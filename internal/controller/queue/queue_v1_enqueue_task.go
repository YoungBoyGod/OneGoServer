package queue

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/queue/v1"
)

func (c *ControllerV1) EnqueueTask(ctx context.Context, req *v1.EnqueueTaskReq) (res *v1.EnqueueTaskRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
