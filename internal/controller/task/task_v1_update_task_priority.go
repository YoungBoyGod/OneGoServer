package task

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/task/v1"
)

func (c *ControllerV1) UpdateTaskPriority(ctx context.Context, req *v1.UpdateTaskPriorityReq) (res *v1.UpdateTaskPriorityRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
