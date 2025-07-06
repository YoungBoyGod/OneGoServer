package queue

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/queue/v1"
)

func (c *ControllerV1) RemoveQueueTask(ctx context.Context, req *v1.RemoveQueueTaskReq) (res *v1.RemoveQueueTaskRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
