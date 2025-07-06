package task

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/task/v1"
)

func (c *ControllerV1) ReassignTask(ctx context.Context, req *v1.ReassignTaskReq) (res *v1.ReassignTaskRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
