package task

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/task/v1"
)

func (c *ControllerV1) AssignTask(ctx context.Context, req *v1.AssignTaskReq) (res *v1.AssignTaskRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
